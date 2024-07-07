package before

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

// AddSeekHelp
// 判断能否新建求助帖子
// 1. 用户是否登录
// 2. 网站若不允许发布求助帖子，该用户是否是管理员
// 3. 该用户是否具有发布求助帖子的权限
// 4. 该用户的悬赏金额是否大于零
func AddSeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reward int
		userId := c.GetInt("id")
		err := models.DB.Model(&models.User{ID: userId}).Select("reward").Scan(&reward).Error
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			c.Abort()
		}
		if reward <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Your balance is insufficient",
			})
			c.Abort()
		}
		c.Set("reward", reward)
		c.Next()
	}
}
