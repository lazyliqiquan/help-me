package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/routes"
	"github.com/lazyliqiquan/help-me/service"
	"log"
	math_rand "math/rand"
	"os"

	"github.com/lazyliqiquan/help-me/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		log.Fatalln("Init logger fail :", err)
	}
	models.Init(logger, config.Config)
	middlewares.Init(logger.Sugar())
	service.Init(logger.Sugar())
	initRand()
	initFiles(config.Config)
	r := routes.Router(config.Config)
	if config.Config.Debug {
		r.Run(config.Config.DebugWebPath)
	} else {
		r.Run(config.Config.WebPath)
	}
}

func initRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		logger.Fatal("Random generator init failed : ", zap.Error(err))
	}
	sd := int64(binary.LittleEndian.Uint64(b[:]))
	logger.Sugar().Infof("random seed : %d ", sd)
	math_rand.Seed(sd)
}

func initFiles(config *config.WebConfig) {
	dirs := []string{config.CodeFilePath, config.ImageFilePath, config.AvatarFilePath}
	for _, v := range dirs {
		if err := os.MkdirAll(v, 0755); err != nil {
			logger.Fatal("Init files create fail : ", zap.Error(err))
		}
	}
}
