package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/service"
	"github.com/lazyliqiquan/help-me/utils"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// Register
// @Tags 公共方法
// @Summary 注册新用户，第一个注册的用户是管理员
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param registerTime formData string true "registerTime"
// @Success 200 {string} json "{"code":"0"}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	registerTime := c.PostForm("registerTime")
	if utils.IsNuiStrs(email, userCode, name, password, registerTime) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code has expired",
			})
		} else {
			service.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failure",
			})
		}
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The verification code is incorrect",
		})
		return
	}
	// 判断邮箱和用户名是否已存在
	var cnt int64
	err = models.DB.Model(&models.User{}).
		Where("email = ? OR name = ?", email, name).Count(&cnt).Error
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The email address or user name already exists",
		})
		return
	}
	var userCount int64
	err = models.DB.Model(&models.User{}).Count(&userCount).Error
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	userBan := 1
	if userCount == 0 {
		// 将第一位注册的用户升级为超级管理员
		userBan = 0
	}
	user := &models.User{
		Name:         name,
		Email:        email,
		Password:     password,
		Reward:       config.Config.UserInitReward,
		RegisterTime: registerTime,
		Ban:          userBan,
	}
	err = models.DB.Create(user).Error
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	// 成功创建用户后，应该立即将验证码销毁掉，以免一把钥匙打开多道们的情况出现
	err = models.RDB.Unlink(c, email).Err()
	if err != nil {
		service.Logger.Errorln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Account registration successful",
	})
}
