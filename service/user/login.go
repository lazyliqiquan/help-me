package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Accept multipart/form-data
// @Param loginType formData string true "loginType"
// @Param nameOrMail formData string true "nameOrMail"
// @Param authCode formData string true "authCode"
// @Success 200 {string} json "{"code":"0"}"
// @Router /login [post]
func Login(c *gin.Context) {
	loginType := c.PostForm("loginType")
	nameOrMail := c.PostForm("nameOrMail")
	authCode := c.PostForm("authCode")
	if utils.IsNuiStrs(loginType, nameOrMail, authCode) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	user := &models.User{}
	// 有三种登录方式：0 : 用户名+密码，1 : 邮箱+密码，2 : 邮箱+验证码
	if loginType == "0" || loginType == "1" {
		condition := "name"
		if loginType == "1" {
			condition = "email"
		}
		condition += " = ? AND password = ?"
		err := models.DB.Model(&models.User{}).
			Where(condition, nameOrMail, authCode).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The user name does not exist or the password is incorrect",
				})
				return
			}
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failure",
			})
			return
		}
	} else {
		sysCode, err := models.RDB.Get(c, nameOrMail).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The mailbox does not exist or the verification code has expired",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failure",
				})
			}
			return
		}
		if sysCode != authCode {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code is incorrect",
			})
			return
		}
		err = models.DB.Model(&models.User{}).
			Where("email = ?", nameOrMail).First(user).Error
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mailbox does not exist",
			})
			return
		}
	}
	// 检查一下用户被封禁没有
	if !models.JudgePermit(models.Login, user.Ban) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The user has been banned",
		})
		return
	}
	// 上面是局部鉴权，下面是全局鉴权
	safeBan, err := models.RDB.Get(c, "safeBan").Result()
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Redis operation failed",
		})
		return
	}
	if safeBan != utils.Permit && !models.JudgePermit(models.Admin, user.Ban) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The website is currently in secure mode, and only administrators can log in",
		})
		return
	}
	token, err := middlewares.GenerateToken(user.ID)
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Generate token fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Login success!",
		"data": gin.H{
			"token":    token,
			"name":     user.Name,
			"password": user.Password,
			"ban":      user.Ban, //前端根据用户权限，创建一些管理员特有的组件
		},
	})
}
