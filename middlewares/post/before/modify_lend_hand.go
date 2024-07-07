package before

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
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
		seekHelpId := c.GetInt("seekHelpId")
		lendHandId := c.GetInt("lendHandId")
		err := models.DB.Model(&models.Post{}).First(&models.Post{}, "id = ?", seekHelpId).Error
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
		err = models.DB.Model(&models.Post{}).Preload("User").Where("id = ?", lendHandId).Select("id", "ban", "status").First(post).Error
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
		userBan := c.GetInt("ban")
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

		c.Next()
	}
}
