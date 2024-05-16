package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"runtime"
)

type WebConfig struct {
	// 系统配置
	SenderMailbox          string `yaml:"senderMailbox"`
	SmtpServerPath         string `yaml:"smtpServerPath"`
	SmtpServerPort         string `yaml:"smtpServerPort"`
	SmtpServerVerification string `yaml:"smtpServerVerification"`
	// 网站配置
	TokenDuration            int `yaml:"tokenDuration"`
	VerificationCodeDuration int `yaml:"verificationCodeDuration"`
	UserInitReward           int `yaml:"userInitReward"`
	MaxDocumentHeight        int `yaml:"maxDocumentHeight"`
	MaxDocumentLength        int `yaml:"maxDocumentLength"`
	MaxPictureSize           int `yaml:"maxPictureSize"`
	MaxCodeFileSize          int `yaml:"maxCodeFileSize"`
	DayUserCommentLimit      int `yaml:"dayUserCommentLimit"`
	TodayShareCodeSurplus    int `yaml:"todayShareCodeSurplus"`
	//权限配置
	SafeBan              string `yaml:"safeBan"`
	PublishSeekHelpBan   string `yaml:"publishSeekHelpBan"`
	ViewSeekHelpBan      string `yaml:"viewSeekHelpBan"`
	LoginViewSeekHelpBan string `yaml:"loginViewSeekHelpBan"`
	ModifySeekHelpBan    string `yaml:"modifySeekHelpBan"`
	PublishLendHandBan   string `yaml:"publishLendHandBan"`
	ViewLendHandBan      string `yaml:"viewLendHandBan"`
	LoginViewLendHandBan string `yaml:"loginViewLendHandBan"`
	ModifyLendHandBan    string `yaml:"modifyLendHandBan"`
	PublishCommentBan    string `yaml:"publishCommentBan"`
	ViewCommentBan       string `yaml:"viewCommentBan"`
	LoginViewCommentBan  string `yaml:"loginViewCommentBan"`
	ModifyCommentBan     string `yaml:"modifyCommentBan"`
	// 不应该修改的配置
	MysqlPassword      string `yaml:"mysqlPassword"`
	WebPath            string `yaml:"webPath"`
	MysqlPath          string `yaml:"mysqlPath"`
	MysqlPort          string `yaml:"mysqlPort"`
	RedisPath          string `yaml:"redisPath"`
	RedisPort          string `yaml:"redisPort"`
	RootFilePath       string `yaml:"rootFilePath"`
	CodeFilePath       string `yaml:"codeFilePath"`
	ImageFilePath      string `yaml:"imageFilePath"`
	AvatarFilePath     string `yaml:"avatarFilePath"`
	MaxIdleConnects    int    `yaml:"maxIdleConnects"`
	MaxMultipartMemory int    `yaml:"maxMultipartMemory"`
	// 是否处于debug
	Debug bool
	// debug
	DebugWebPath   string `yaml:"debugWebPath"`
	DebugMysqlPath string `yaml:"debugMysqlPath"`
	DebugMysqlPort string `yaml:"debugMysqlPort"`
	DebugRedisPath string `yaml:"debugRedisPath"`
	DebugRedisPort string `yaml:"debugRedisPort"`
}

func (c *WebConfig) RedisInit() map[string]any {
	res := make(map[string]any)
	// 系统配置
	res["senderMailbox"] = c.SenderMailbox
	res["smtpServerPath"] = c.SmtpServerPath
	res["smtpServerPort"] = c.SmtpServerPort
	res["smtpServerVerification"] = c.SmtpServerVerification
	// 网站配置
	//res["tokenDuration"] = c.TokenDuration
	//res["verificationCodeDuration"] = c.VerificationCodeDuration
	//res["userInitScore"] = c.UserInitReward
	res["maxDocumentHeight"] = c.MaxDocumentHeight
	res["maxDocumentLength"] = c.MaxDocumentLength
	res["maxPictureSize"] = c.MaxPictureSize
	res["maxCodeFileSize"] = c.MaxCodeFileSize
	res["dayUserCommentLimit"] = c.DayUserCommentLimit
	res["todayShareCodeSurplus"] = c.TodayShareCodeSurplus
	// 权限配置
	res["safeBan"] = c.SafeBan
	res["publishSeekHelpBan"] = c.PublishSeekHelpBan
	res["viewSeekHelpBan"] = c.ViewSeekHelpBan
	res["loginViewSeekHelpBan"] = c.LoginViewSeekHelpBan
	res["modifySeekHelpBan"] = c.ModifySeekHelpBan
	res["publishLendHandBan"] = c.PublishLendHandBan
	res["viewLendHandBan"] = c.ViewLendHandBan
	res["loginViewLendHandBan"] = c.LoginViewLendHandBan
	res["modifyLendHandBan"] = c.ModifyLendHandBan
	res["publishCommentBan"] = c.PublishCommentBan
	res["viewCommentBan"] = c.ViewCommentBan
	res["loginViewCommentBan"] = c.LoginViewCommentBan
	res["modifyCommentBan"] = c.ModifyCommentBan
	return res
}

var Config *WebConfig

func init() {
	Config = &WebConfig{}
	if runtime.GOOS == "windows" {
		Config.Debug = true
	}
	file, err := os.Open("settings.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	//defer file.Close()//一定要处理关闭错误吗？
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(data, Config)
	if err != nil {
		log.Fatalln(err)
	}
	Config.CodeFilePath = Config.RootFilePath + Config.CodeFilePath
	Config.ImageFilePath = Config.RootFilePath + Config.ImageFilePath
	Config.AvatarFilePath = Config.RootFilePath + Config.AvatarFilePath
}
