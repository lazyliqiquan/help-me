package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

type OneCommentBan struct {
	CommentId int   `form:"commentId" binding:"required"`
	Ban       *bool `form:"ban" binding:"required"`
}

// ForbidOneComment 操作某条评论的权限
// @Tags 管理员方法
// @Summary 封禁某条评论
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param commentId formData int true "1"
// @Param ban formData bool true "false"
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/forbid-one-comment [post]
func ForbidOneComment(c *gin.Context) {
	var oneCommentBan OneCommentBan
	if err := c.ShouldBind(&oneCommentBan); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	err := models.DB.Model(&models.Comment{ID: oneCommentBan.CommentId}).Update("ban", *oneCommentBan.Ban).Error
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
		"msg":  "The comments were successfully blocked",
	})
}
