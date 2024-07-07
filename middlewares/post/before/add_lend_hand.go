package before

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

// AddLendHand
// 1. 该用户是否登录
// 2. 若网站不允许发布帮助帖子，该用户是否是管理员
// 3. 该用户是否具有发布帮助帖子的权限
// 4. 传递到后端的求助帖子ID是否存在(url可能是用户自行键入的，有错误的情况)
// 5. 该用户此前是否已经针对该求助帖子发布过帮助帖子
// 6. 若求助帖子不允许发布新的帮助帖子，该用户是否是管理员
func AddLendHand() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("id")
		post := &models.Post{}
		post.ID = c.GetInt("seekHelpId")
		err := models.DB.Model(&models.Post{}).Preload("LendHands", "user_id = ?", userId).
			Where("id = ?", post.ID).Select("id", "ban").First(post).Error
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
		if len(post.LendHands) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Only one lend hand post can be posted for each seek help post",
			})
			c.Abort()
		}
		userBan := c.GetInt("ban")
		if !models.JudgePermit(models.AddLendHand, post.Ban) && !models.JudgePermit(models.Admin, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "This seek help post cannot be assisted and permissions are disabled",
			})
			c.Abort()
		}
		c.Next()
	}
}
