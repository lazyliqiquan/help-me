package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

// Add
// 1. 该用户是否登录
// 2. 若网站不允许发布评论，该用户是否是管理员
// 3. 该用户是否具有发布评论的权限
// 4. 传递到后端的帖子id是否存在
// 5. 若帖子不允许发布新的评论，该用户是否是管理员
// 6. 用户发布的评论是否达到每日上限
// @Tags 用户方法
// @Summary 添加新的评论
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postId formData int true "1"
// @Param text formData string true "Hello"
// @Param sendTime formData string true "2024-05-26 15:10:00"
// @Param isBefore formData string true "0"
// @Success 200 {string} json "{"code":"0"}"
// @Router /add-comment [post]
func Add(c *gin.Context) {
	postId := c.GetInt("postId")
	userId := c.GetInt("userId")
	commentSum := c.GetInt("commentSum")
	text := c.GetString("text")
	time := c.PostForm("time")
	comment := &models.Comment{
		Time: time,
		Text: text,
		//这里直接添加，关联模式会自动维护吗
		PostID: postId,
		UserID: userId,
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Comment{}).Create(comment).Error
		if err != nil {
			return err
		}
		return tx.Model(&models.Post{ID: postId}).Update("comment_sum", commentSum+1).
			Association("Comments").Append(&models.Comment{ID: comment.ID})
	})
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
		"msg":  "Add a comment successfully",
	})
}
