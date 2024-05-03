package before

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/service"
	"net/http"
)

// AddSeekHelp
// @Tags 用户方法
// @Summary 检查该用户是否具有新建求助帖子的资格
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Success 200 {string} json "{"code":"0"}"
// @Router /before/add-seek-help [post]
func AddSeekHelp(c *gin.Context) {
	userId := c.GetInt("id")
	var reward int
	if err := models.DB.Model(&models.User{}).Where("id = ?", userId).Select("reward").Scan(&reward).Error; err != nil {
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failed",
		})
		return
	}
	if reward <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "You have no amount to issue a reward",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "You are ready to start editing",
	})
}
