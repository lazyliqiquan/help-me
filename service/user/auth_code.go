package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/service"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"time"
)

// SendCode
// @Tags 公共方法
// @Summary 发送验证码(一个验证码只能处理一个操作，用完就要删除)
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"0"}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The mailbox cannot be empty",
		})
		return
	}
	// _, err := models.RDB.Get(c, email).Result()
	ttlResult, err := models.RDB.TTL(c, email).Result()
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The redis operation failed",
		})
		return
	} else if ttlResult == time.Duration(-1) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "No expiration time is set for the current Key",
		})
		return
	} else if ttlResult != time.Duration(-2) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The verification code has not expired, please use the previous one",
			"data": gin.H{
				"ttl": ttlResult.Seconds(),
			},
		})
		return
	}
	code := utils.GetRand()
	err = utils.SendCode(email, code)
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Failed to send the verification code",
		})
		return
	}
	err = models.RDB.Set(c, email, code,
		time.Duration(config.Config.VerificationCodeDuration*int(time.Minute))).Err()
	if err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Unable to write data to redis",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "The verification code is sent successfully",
		"data": gin.H{
			"ttl": time.Duration(config.Config.VerificationCodeDuration * int(time.Minute)).Seconds(),
		},
	})
}
