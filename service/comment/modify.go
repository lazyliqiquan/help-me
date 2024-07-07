package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

type NewComment struct {
	Id   int    `form:"id" binding:"required"`
	Text string `form:"text" binding:"required"`
	Time string `form:"time" binding:"required"`
}

// Modify
// 1. 该用户是否登录
// 2. 若网站不允许修改评论，该用户是否是管理员
// 3. 该用户是否具有修改评论的权限
// 4. 该用户是否是该评论的所有者
// @Tags 用户方法
// @Summary 修改评论
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param id formData int true "1"
// @Param text formData string true "Hello"
// @Param time formData string true "2024-05-26 15:10:00"
// @Success 200 {string} json "{"code":"0"}"
// @Router /modify-comment [post]
func Modify(c *gin.Context) {
	var newComment NewComment
	if err := c.ShouldBind(&newComment); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	comment := &models.Comment{}
	err := models.DB.Model(&models.Comment{ID: newComment.Id}).First(comment).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	userId := c.GetInt("id")
	userBan := c.GetInt("ban")
	//管理员可以修改所有的评论
	if !models.JudgePermit(models.Admin, userBan) {
		//该用户不是该条评论的所有者，无法修改
		if comment.UserID != userId {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not the owner of this comment and cannot modify it",
			})
			return
		}
	}
	comment.Time = newComment.Time
	comment.Text = newComment.Text
	err = models.DB.Save(comment).Error
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
		"msg":  "Modified comment successfully",
	})
}
