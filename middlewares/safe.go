package middlewares

import (
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
)

// TokenSafeModel 安全模式
// 附带解析token(后续GetInt("id")为0，表明用户没有有效token或者处于封禁状态)
func TokenSafeModel() gin.HandlerFunc {
	return func(c *gin.Context) {
		safeBan, err := models.RDB.Get(c, "safeBan").Result()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		var userClaim *UserClaims
		var userClaimErr error
		var userBan int
		auth := c.GetHeader("Authorization")
		if auth != "" {
			userClaim, userClaimErr = AnalyseToken(auth)
			if userClaimErr != nil {
				utils.Logger.Errorln(userClaimErr)
			} else {
				// 向数据库获取最新的user ban
				err := models.DB.Model(&models.User{ID: userClaim.Id}).
					Select("ban").Scan(&userBan).Error
				if err != nil {
					utils.Logger.Errorln(err)
					c.JSON(http.StatusOK, gin.H{
						"code": 1,
						"msg":  "Mysql operation failed",
					})
					c.Abort()
				}
				// 局部检查该用户是否被封禁
				// 虽然login那里就已经判断了，但是可能会出现用户登陆后，管理员在token过期之前进行封禁的情况
				if !models.JudgePermit(models.Login, userBan) {
					c.JSON(http.StatusOK, gin.H{
						"code": 1,
						"msg":  "The user has been banned",
					})
					c.Abort()
				}
				c.Set("id", userClaim.Id)
				c.Set("ban", userBan)
			}
		}
		if safeBan != utils.Permit && (auth == "" || userClaimErr != nil || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		utils.Logger.Infoln("TokenSafeModel")
		c.Next()
	}
}

// OtherSafeModel 主要针对找回密码和注册
func OtherSafeModel() gin.HandlerFunc {
	return func(c *gin.Context) {
		safeBan, err := models.RDB.Get(c, "safeBan").Result()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		if safeBan != utils.Permit {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The site is currently in safe mode and this feature is not available",
			})
			c.Abort()
		}
		utils.Logger.Infoln("OtherSafeModel")
		c.Next()
	}
}

// LoginModel 判断是否登录
func LoginModel() gin.HandlerFunc {
	return func(c *gin.Context) {
		// id不可能为0，如果为0，说明用户未登录
		userId := c.GetInt("id")
		if userId == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Login is required for this operation",
			})
			c.Abort()
		}
		utils.Logger.Infoln("LoginModel")
		c.Next()
	}
}
