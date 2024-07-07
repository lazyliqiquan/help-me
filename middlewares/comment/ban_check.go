package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

var (
	commentBanList = []string{"modifyCommentBan", "publishCommentBan"}
	userBanList    = []int{models.PublishComment, models.ModifyComment}
)

// BanCheck
// 评论权限检查
func BanCheck(isAdd int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userBan := c.GetInt("ban")
		if !models.JudgePermit(models.Admin, userBan) {
			modifyBan, err := models.RDB.Get(c, commentBanList[isAdd]).Result()
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failed",
				})
				c.Abort()
			}
			//全局判断
			if modifyBan != utils.Permit {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The website is currently in" + commentBanList[isAdd] + " mode and can only be operated by administrators",
				})
				c.Abort()
			}
			//个人判断
			if !models.JudgePermit(userBanList[isAdd], userBan) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Your permission is blocked",
				})
				c.Abort()
			}
		}
		utils.Logger.Infoln("comment ban check")
		c.Next()
	}
}
