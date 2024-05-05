package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

var (
	modifyBan     = []string{"modifySeekHelpBan", "modifyLendHandBan", "modifyCommentBan"}
	userModifyBan = []int{models.ModifySeekHelp, models.ModifyLendHand, models.ModifyComment}
)

// Modify
// 预处理，判断用户是否具有修改权限
func Modify(modifyType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userBan := c.GetInt("ban")
		if !models.JudgePermit(models.Admin, userBan) {
			modifyBan, err := models.RDB.Get(c, modifyBan[modifyType]).Result()
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
					"msg":  "The website is currently in modify safe mode and can only be operated by administrators",
				})
				c.Abort()
			}
			//个人判断
			if !models.JudgePermit(userModifyBan[modifyType], userBan) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Your permission to modify is blocked",
				})
				c.Abort()
			}
		}
		utils.Logger.Infoln("Modify judge")
		c.Next()
	}
}
