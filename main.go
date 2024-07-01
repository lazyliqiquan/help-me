package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/routes"
	"github.com/lazyliqiquan/help-me/utils"
	math_rand "math/rand"
	"os"
	"strconv"

	"github.com/lazyliqiquan/help-me/config"
	_ "github.com/lazyliqiquan/help-me/docs"
	"go.uber.org/zap"
)

func main() {
	initRand()
	initFiles(config.Config)
	models.Init(config.Config)
	r := routes.Router(config.Config)
	webPath := config.Config.WebPath
	if config.Config.Debug {
		webPath = config.Config.DebugWebPath
	}
	if err := r.Run(webPath); err != nil {
		utils.RootLogger.Fatal("Router run fail : ", zap.Error(err))
	}
}

func initRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		utils.RootLogger.Fatal("Random generator init failed : ", zap.Error(err))
	}
	sd := int64(binary.LittleEndian.Uint64(b[:]))
	utils.Logger.Infof("random seed : %d ", sd)
	math_rand.Seed(sd)
}

func initFiles(config *config.WebConfig) {
	dirs := []string{config.ImageFilePath, config.AvatarFilePath}
	for _, v := range dirs {
		if err := os.MkdirAll(v, 0755); err != nil {
			utils.RootLogger.Fatal("Init root dir create fail : ", zap.Error(err))
		}
	}
	for i := 0; i < config.SecondImageDirAmount; i++ {
		secondDirName := config.ImageFilePath + strconv.Itoa(i) + "/"
		if err := os.MkdirAll(secondDirName, 0755); err != nil {
			utils.RootLogger.Fatal("Init second dir create fail : ", zap.Error(err))
		}
	}
}
