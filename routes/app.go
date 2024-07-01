package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"time"
)

func Router(config *config.WebConfig) *gin.Engine {
	r := gin.Default()
	if config.Debug {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Access-Control-Allow-Origin"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	r.MaxMultipartMemory = int64(config.MaxMultipartMemory) << 20
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes(r)
	if !config.Debug {
		err := r.RunTLS(config.WebPath, "assets/https/cert.pem", "assets/https/key.pem")
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}
