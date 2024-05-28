package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/middlewares/post/before"
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
	postType := c.PostForm("postType")
	restrictions := make(map[string]int)
	//fixme 在middleware里面执行c.Next()和c.Abort()会发生意料之外的事情吗？
	switch postType {
	case "0":
		userId := c.GetInt("id")
		var reward int
		err := models.DB.Model(&models.User{ID: userId}).Select("reward").Scan(&reward).Error
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			return
		}
		restrictions["reward"] = reward
	case "1":
		before.AddLendHand()(c)
	case "2":
		before.ModifySeekHelp()(c)
	case "3":
		before.ModifyLendHand()(c)
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Post type does not match any of the cases",
		})
		return
	}
	for k := range config.Config.GetRestrictionSetting() {
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
