package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares/post/before"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"strconv"
)

var (
	postBanList = []string{"publishSeekHelpBan", "publishLendHandBan", "modifySeekHelpBan", "modifyLendHandBan"}
	userBanList = []int{models.PublishSeekHelp, models.PublishLendHand, models.ModifySeekHelp, models.ModifyLendHand}
)

// BanCheck
// 在进行有关帖子的操作之前的检查和处理
func BanCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userBan := c.GetInt("ban")
		postType, err := strconv.Atoi(c.PostForm("postType"))
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Post type param is not integer",
			})
			c.Abort()
		}
		if !models.JudgePermit(models.Admin, userBan) {
			globalPostBan, err := models.RDB.Get(c, postBanList[postType]).Result()
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failed",
				})
				c.Abort()
			}
			//全局判断
			if globalPostBan != utils.Permit {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The website is currently in " + postBanList[postType] + " mode and can only be operated by administrators",
				})
				c.Abort()
			}
			//个人判断
			if !models.JudgePermit(userBanList[postType], userBan) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Your permission is blocked",
				})
				c.Abort()
			}
		}
		//	====================上面是全局检查，下面是局部检查======================
		if postType > 0 {
			seekHelpId, err := strconv.Atoi(c.PostForm("seekHelpId"))
			if err != nil {
				utils.Logger.Errorln(err)
				// 不是整数的情况应该交给前端处理，我们不需要额外说明
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Seek help id nonentity",
				})
				c.Abort()
			}
			c.Set("seekHelpId", seekHelpId)
		}
		if postType == 3 {
			lendHandId, err := strconv.Atoi(c.PostForm("lendHandId"))
			if err != nil {
				utils.Logger.Errorln(err)
				// 不是整数的情况应该交给前端处理，我们不需要额外说明
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Lend hand id nonentity",
				})
				c.Abort()
			}
			c.Set("lendHandId", lendHandId)
		}
		c.Set("postType", postType)
		//	中间件里面套中间件，应该没事吧
		switch postType {
		case 0:
			before.AddSeekHelp()(c)
		case 1:
			before.AddLendHand()(c)
		case 2:
			before.ModifySeekHelp()(c)
		case 3:
			before.ModifyLendHand()(c)
		default:
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Post type does not match any of the cases",
			})
			c.Abort()
		}
		utils.Logger.Infoln("post ban check")
		c.Next()
	}

}
