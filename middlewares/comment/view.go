package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

const (
	viewBan      = "viewCommentBan"
	loginViewBan = "loginViewCommentBan"
)

// View
// 预处理，判断用户是否具有查看权限
func View() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewBan, err := models.RDB.Get(c, viewBan).Result()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		loginViewBan, err := models.RDB.Get(c, loginViewBan).Result()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		//如果userId为0表示未登录
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		if viewBan != utils.Permit && (userId == 0 || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in view safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		if loginViewBan != utils.Permit && userId == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not logged in and cannot browse",
			})
			c.Abort()
		}
		utils.Logger.Infoln("Comment view")
		c.Next()
	}
}
