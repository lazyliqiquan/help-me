package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

type ForbidParam struct {
	PostId int `form:"postId" binding:"required"`
	//修改的是那种权限
	BanOption int `form:"banOption" binding:"required"`
	//true 给帖子添加某种权限 false 给帖子撤销某种权限
	IsAdd bool `form:"isAdd" binding:"required"`
}

// ForbidPost 操作帖子权限
// @Tags 管理员方法
// @Summary 修改帖子权限
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postId formData int true 1
// @Param banOption formData int true 0
// @Param isAdd formData bool true false
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/forbid-post [post]
func ForbidPost(c *gin.Context) {
	var forbidParam ForbidParam
	if err := c.ShouldBind(&forbidParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	var postBan int
	err := models.DB.Model(&models.Post{ID: forbidParam.PostId}).Select("ban").Scan(&postBan).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	if forbidParam.IsAdd {
		postBan = models.AddPermit(forbidParam.BanOption, postBan)
	} else {
		postBan = models.SubPermit(forbidParam.BanOption, postBan)
	}
	err = models.DB.Model(&models.Post{ID: forbidParam.PostId}).Update("ban", postBan).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Succeeded in modifying the post permission",
	})
}
