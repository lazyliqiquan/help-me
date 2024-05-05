package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

var (
	publishBan     = []string{"publishSeekHelpBan", "publishLendHandBan", "publishCommentBan"}
	userPublishBan = []int{models.PublishSeekHelp, models.PublishLendHand, models.PublishComment}
)

// Publish
// 预处理，判断用户是否具有修改权限
func Publish(publishType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userBan := c.GetInt("ban")
		if !models.JudgePermit(models.Admin, userBan) {
			publishBan, err := models.RDB.Get(c, publishBan[publishType]).Result()
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failed",
				})
				c.Abort()
			}
			if publishBan != utils.Permit {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The website is currently in publish safe mode and can only be operated by administrators",
				})
				c.Abort()
			}
			if !models.JudgePermit(userPublishBan[publishType], userBan) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Your permission to publish is blocked",
				})
				c.Abort()
			}
		}
		utils.Logger.Infoln("Publish judge")
		c.Next()
	}
}
