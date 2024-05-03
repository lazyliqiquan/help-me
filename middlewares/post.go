package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
)

// 浏览帖子
func ViewPost(isSeekHelp bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		viewPostBanStr := "viewSeekHelpBan"
		loginViewPostBanStr := "loginViewSeekHelpBan"
		if !isSeekHelp {
			viewPostBanStr = "viewLendHandBan"
			loginViewPostBanStr = "loginViewLendHandBan"
		}
		viewPostBan, err := models.RDB.Get(c, viewPostBanStr).Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		loginViewPostBan, err := models.RDB.Get(c, loginViewPostBanStr).Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		// 全局判断
		if viewPostBan != "permit" && (userId == 0 || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		if loginViewPostBan != "permit" && (userId == 0 || !models.JudgePermit(models.Login, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not logged in or do not have browsing rights",
			})
			c.Abort()
		}
		logger.Infoln("ViewSeekHelp")
		c.Next()
	}
}
