package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
)

var (
	editBan     = []string{"publishSeekHelpBan", "editSeekHelpBan", "publishLendHandBan", "editLendHandBan"}
	userEditBan = []int{models.PublishSeekHelp, models.EditSeekHelp, models.PublishLendHand, models.EditLendHand}
)

// 判断用户是否具有编辑权限(新增求助,修改求助,新增帮助,修改帮助)
func JudgeEdit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userBan := c.GetInt("ban")
		editOption, err := strconv.Atoi(c.PostForm("editOption"))
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Param type not is a int type",
			})
			c.Abort()
		}
		// 如果是管理员就不需要进行判断
		if !models.JudgePermit(models.Admin, userBan) {
			isEditBan, err := models.RDB.Get(c, editBan[editOption]).Result()
			if err != nil {
				logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failed",
				})
				c.Abort()
			}
			// 全局判断
			if isEditBan != "permit" {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "This feature is disabled on the website",
				})
				c.Abort()
			}
			// 局部判断(这里只是判断用户的权限，但是还没有判断该帖子的权限)
			if !models.JudgePermit(userEditBan[editOption], userBan) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "You do not have permission to do this",
				})
				c.Abort()
			}
		}
		// 因为c.GetInt如果获取不到数据，会返回0，所以为了后面区分信息，这里还是使用正整数好一点
		c.Set("editOption", editOption)
		logger.Infoln("JudgeEdit")
		c.Next()
	}
}
