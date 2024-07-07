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
	MaxCommentWords     int `yaml:"maxCommentWords"`
	MaxDocumentWords    int `yaml:"maxDocumentWords"`
	MaxPicturesSize     int `yaml:"maxPicturesSize"`
	DayUserCommentLimit int `yaml:"dayUserCommentLimit"`
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
	TokenDuration            int    `yaml:"tokenDuration"`
	VerificationCodeDuration int    `yaml:"verificationCodeDuration"`
	UserInitReward           int    `yaml:"userInitReward"`
	MysqlPassword            string `yaml:"mysqlPassword"`
	WebPath                  string `yaml:"webPath"`
	MysqlPath                string `yaml:"mysqlPath"`
	MysqlPort                string `yaml:"mysqlPort"`
	RedisPath                string `yaml:"redisPath"`
	RedisPort                string `yaml:"redisPort"`
	RootFilePath             string `yaml:"rootFilePath"`
	ImageFilePath            string `yaml:"imageFilePath"`
	AvatarFilePath           string `yaml:"avatarFilePath"`
	SecondImageDirAmount     int    `yaml:"secondImageDirAmount"`
	MaxIdleConnects          int    `yaml:"maxIdleConnects"`
	MaxMultipartMemory       int    `yaml:"maxMultipartMemory"`
	// 是否处于debug
	Debug bool
	// debug
	DebugWebPath   string `yaml:"debugWebPath"`
	DebugMysqlPath string `yaml:"debugMysqlPath"`
	DebugMysqlPort string `yaml:"debugMysqlPort"`
	DebugRedisPath string `yaml:"debugRedisPath"`
	DebugRedisPort string `yaml:"debugRedisPort"`
}

func (c *WebConfig) GetRestrictionSetting() map[string]int {
	restriction := make(map[string]int)
	// 网站配置
	restriction["maxCommentWords"] = c.MaxCommentWords
	restriction["maxDocumentWords"] = c.MaxDocumentWords
	restriction["maxPicturesSize"] = c.MaxPicturesSize
	restriction["dayUserCommentLimit"] = c.DayUserCommentLimit
	return restriction
}

func (c *WebConfig) GetPermissionSetting() map[string]string {
	permissionSetting := make(map[string]string)
	// 权限配置
	permissionSetting["safeBan"] = c.SafeBan
	permissionSetting["publishSeekHelpBan"] = c.PublishSeekHelpBan
	permissionSetting["viewSeekHelpBan"] = c.ViewSeekHelpBan
	permissionSetting["loginViewSeekHelpBan"] = c.LoginViewSeekHelpBan
	permissionSetting["modifySeekHelpBan"] = c.ModifySeekHelpBan
	permissionSetting["publishLendHandBan"] = c.PublishLendHandBan
	permissionSetting["viewLendHandBan"] = c.ViewLendHandBan
	permissionSetting["loginViewLendHandBan"] = c.LoginViewLendHandBan
	permissionSetting["modifyLendHandBan"] = c.ModifyLendHandBan
	permissionSetting["publishCommentBan"] = c.PublishCommentBan
	permissionSetting["viewCommentBan"] = c.ViewCommentBan
	permissionSetting["loginViewCommentBan"] = c.LoginViewCommentBan
	permissionSetting["modifyCommentBan"] = c.ModifyCommentBan
	return permissionSetting
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
	Config.ImageFilePath = Config.RootFilePath + Config.ImageFilePath
	Config.AvatarFilePath = Config.RootFilePath + Config.AvatarFilePath
}
