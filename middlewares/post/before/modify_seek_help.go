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

// ModifySeekHelp
// 1. 用户是否登录
// 2. 网站若不允许修改求助帖子，该用户是否是管理员
// 3. 该用户是否具有修改求助帖子的权限
// 4. 传递到后端的求助帖子ID是否存在(url可能是用户自行键入的，有错误的情况)
// 5. 若用户不是该求助帖子的拥有者，该用户是否是管理员
// 6. 若该求助帖子不允许修改，该用户是否是管理员
// 7. 若该求助帖子下已经存在帮助帖子，该用户是否是管理员
func ModifySeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		post := &models.Post{}
		var err error
		post.ID, err = strconv.Atoi(c.PostForm("seekHelpId"))
		if err != nil {
			utils.Logger.Errorln(err)
			// 不是整数的情况应该交给前端处理，我们不需要额外说明
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Seek help id nonentity",
			})
			c.Abort()
		}
		//fixme User表明默认加上 "s" ?
		//好像必须要有id才能预加载成功
		err = models.DB.Model(&models.Post{}).Preload("Users").Where("id = ?", post.ID).Select("id", "ban", "lend_hand_sum").First(post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { //传递过来的seekHelpId不存在，用户输入的url有问题
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
		if !models.JudgePermit(models.Admin, userBan) {
			if userId != post.User.ID {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "You are not the owner of the post and cannot modify it",
				})
				c.Abort()
			}
			if !models.JudgePermit(models.Modify, post.Ban) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The current seek help post cannot be modified",
				})
				c.Abort()
			}
			if post.LendHandSum > 0 {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The current seek help post has a sponsor and cannot be modified",
				})
				c.Abort()
			}
		}
		c.Set("seekHelpId", post.ID)
		c.Next()
	}
}
