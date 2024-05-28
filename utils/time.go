package utils

import (
	"go.uber.org/zap"
	"strings"
	"time"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		RootLogger.Fatal("Unable to load time zone ", zap.Error(err))
	}
}

// GetLocationTime 获取北京地区的时间
func GetLocationTime() string {
	return time.Now().In(location).Format("2024-01-02 15:04:05")
}

// GetCurrentDate 获取日期
func GetCurrentDate() string {
	return strings.Split(GetLocationTime(), " ")[0]
}

// GetCurrentTime 获取时间
func GetCurrentTime() string {
	return strings.Split(GetLocationTime(), " ")[1]
}
