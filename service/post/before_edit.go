package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

// BeforeEdit 编辑操作前的预判断
// @Tags 用户方法
// @Summary 编辑操作前的预判断
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postType formData string true "0"
// @Param seekHelpId formData int false "1"
// @Param lendHandId formData int false "2"
// @Success 200 {string} json "{"code":"0"}"
// @Router /before-edit [post]
func BeforeEdit(c *gin.Context) {
	postType := c.GetInt("postType")
	restrictions := make(map[string]int)
	if postType == 0 {
		restrictions["reward"] = c.GetInt("reward")
	}
	for _, k := range []string{"maxDocumentWords", "maxPicturesSize"} {
		v, err := models.RDB.Get(c, k).Int()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis error",
			})
			return
		}
		restrictions[k] = v
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "You gain access to the edit screen",
		"data": gin.H{"map": restrictions},
	})
}
