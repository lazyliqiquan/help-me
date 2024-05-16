package before

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ModifyLendHand
// 1. 该用户是否登录
// 2. 网站若不允许修改帮助帖子，该用户是否是管理员
// 3. 该用户是否具有修改帮助帖子的权限
// 4. 传递到后端的求助帖子ID是否存在(url可能是用户自行键入的，有错误的情况)
// 5. 传递到后端的帮助帖子ID是否存在(url可能是用户自行键入的，有错误的情况)
// 6. 若该帮助帖子不允许修改，该用户是否是管理员
// 7. 若该帮助帖子已经被求助者接受，该用户是否是管理员
// 8. 若该用户不是该帮助帖子的拥有者，该用户是否是管理员
func ModifyLendHand() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		seekHelpId, err := strconv.Atoi(c.PostForm("seekHelpId"))
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Seek help id nonentity",
			})
			c.Abort()
		}
		lendHandId, err := strconv.Atoi(c.PostForm("lendHandId"))
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Lend hand id nonentity",
			})
			c.Abort()
		}
		err = models.DB.Model(&models.Post{}).First(&models.Post{}, "id = ?", seekHelpId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Url error : seek help id not exist",
				})
				c.Abort()
			}
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			c.Abort()
		}
		post := &models.Post{}
		//fixme 这里到底是 User 还是 Users
		err = models.DB.Model(&models.Post{}).Preload("Users").Where("id = ?", lendHandId).Select("id", "ban", "status").First(post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Url error : lend hand id not exist",
				})
				c.Abort()
			}
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			c.Abort()
		}
		if !models.JudgePermit(models.Admin, userBan) {
			if !models.JudgePermit(models.Modify, post.Ban) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The current lend hand post cannot be modified",
				})
				c.Abort()
			}
			if post.Status {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The current help post has been accepted and cannot be modified",
				})
				c.Abort()
			}
			if userId != post.User.ID {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "You are not the owner of the post and cannot modify it",
				})
				c.Abort()
			}
		}
		c.Set("seekHelpId", seekHelpId)
		c.Set("lendHandId", lendHandId)
		c.Next()
	}
}
