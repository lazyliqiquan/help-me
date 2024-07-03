package models

import (
	"context"
	"github.com/lazyliqiquan/help-me/utils"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/lazyliqiquan/help-me/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func Init(config *config.WebConfig) {
	var err error
	mysqlDsn := "root:" + config.MysqlPassword + "@tcp("
	if config.Debug {
		mysqlDsn +=
			config.DebugMysqlPath + ":" + config.DebugMysqlPort
	} else {
		mysqlDsn +=
			config.MysqlPath + ":" + config.MysqlPort
	}
	mysqlDsn += ")/?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	if err != nil {
		utils.RootLogger.Fatal("Connect mysql fail : ", zap.Error(err))
	}
	sqlDB, err := DB.DB()
	if err != nil {
		utils.RootLogger.Fatal("Give sql.DB fail : ", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConnects)
	err = DB.Exec("CREATE DATABASE IF NOT EXISTS help_me").Error
	if err != nil {
		utils.RootLogger.Fatal("Create database help_me fail : ", zap.Error(err))
	}
	err = DB.Exec("USE help_me").Error
	if err != nil {
		utils.RootLogger.Fatal("Unable to use the database help_me : ", zap.Error(err))
	}
	err = DB.AutoMigrate(&User{}, &Post{}, &PostStats{}, &Comment{})
	if err != nil {
		utils.RootLogger.Fatal("Create tables fail : ", zap.Error(err))
	}
	utils.Logger.Infoln("Help-me mysql init succeed !")
	redisDsn := config.RedisPath + ":" + config.RedisPort
	if config.Debug {
		redisDsn = config.DebugRedisPath + ":" + config.DebugRedisPort
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisDsn,
		Password: "",
	})
	for k, v := range config.GetRestrictionSetting() {
		if err := RDB.Set(context.Background(), k, v, time.Duration(0)).Err(); err != nil {
			utils.RootLogger.Fatal("Set website setting fail : ", zap.Error(err))
		}
	}
	for k, v := range config.GetPermissionSetting() {
		if err := RDB.Set(context.Background(), k, v, time.Duration(0)).Err(); err != nil {
			utils.RootLogger.Fatal("Set permission setting fail : ", zap.Error(err))
		}
	}
	utils.Logger.Infoln("Help-me redis init succeed !")
	// 启动一个协程来每天重置网站配置
	//go webTicker()
}

func webTicker() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		utils.RootLogger.Fatal("Unable to load time zone ", zap.Error(err))
	}
	now := time.Now().In(location)
	// 格式化时间
	// timeStr := now.Format("2006-01-02 15:04:05")
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nextMidnight = nextMidnight.Add(time.Hour * 24)
	duration := nextMidnight.Sub(now)
	time.Sleep(duration)
	resetWebConfig()
	// 创建定时器，在每天的 0 点触发更新操作
	duration = time.Duration(time.Hour * 24)
	ticker := time.NewTicker(duration)
	for range ticker.C {
		resetWebConfig()
	}
}

// 定期重置网站配置
func resetWebConfig() {
	keys := []string{"daySeekHelpLimit", "dayLendHandLimit", "dayShareCodeLimit"}
	result, err := RDB.MGet(context.Background(), keys...).Result()
	if err != nil {
		utils.Logger.Errorln(err)
		return
	}
	m := map[string]any{
		"todaySeekHelpSurplus":  result[0],
		"todayLendHandSurplus":  result[1],
		"todayShareCodeSurplus": result[2],
	}
	err = RDB.MSet(context.Background(), m).Err()
	if err != nil {
		utils.Logger.Errorln(err)
	}
}
