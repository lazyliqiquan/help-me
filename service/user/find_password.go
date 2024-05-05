package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

// FindPassword
// @Tags 公共方法
// @Summary 找回密码
// @Accept multipart/form-data
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"0"}"
// @Router /find-password [post]
func FindPassword(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	password := c.PostForm("password")
	if utils.IsNuiStrs(email, userCode, password) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code has expired",
			})
		} else {
			utils.Logger.Errorln(err)
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
	user := &models.User{}
	err = models.DB.Model(&models.User{}).Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Email not registered",
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
	err = models.DB.Model(&models.User{}).Where(&models.User{ID: user.ID}).
		Update("password", password).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	err = models.RDB.Unlink(c, email).Err()
	if err != nil {
		utils.Logger.Errorln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Password changed successfully",
	})
}
