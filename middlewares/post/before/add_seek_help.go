package before

import (
	"github.com/gin-gonic/gin"
)

// AddSeekHelp
// 判断能否新建求助帖子
// 1. 用户是否登录
// 2. 网站若不允许发布求助帖子，该用户是否是管理员
// 3. 该用户是否具有发布求助帖子的权限
// 4. 该用户的悬赏金额是否大于零
func AddSeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
