package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/service/user"
)

// 任何人都可以访问
func logout(r *gin.Engine) {
	r.POST("/login", user.Login)
	anyoneAuth := r.Group("/", middlewares.OtherSafeModel())
	anyoneAuth.POST("/send-code", user.SendCode)
	anyoneAuth.POST("/register", user.Register)
	anyoneAuth.POST("/find-password", user.FindPassword)
}
