package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ViewPost
// @Tags 公共方法
// @Summary 浏览求助帖子，但其实修改一下swagger的Router，也可以浏览帮助帖子
// @Accept multipart/form-data
// @Param postId formData string true "1"
// @Success 200 {string} json "{"code":"0"}"
// @Router /view/seek-help [post]
func ViewPost(c *gin.Context) {
	userBan := c.GetInt("ban")
	postId, err := strconv.Atoi(c.PostForm("postId"))
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Post id not a valid integer",
		})
		return
	}
	post := &models.Post{ID: postId}
	err = models.DB.Model(post).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "avatar", "reward", "register_time")
	}).Preload("PostStats").First(post).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	if !models.JudgePermit(models.View, post.Ban) && !models.JudgePermit(models.Admin, userBan) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "You do not have permission to view this post",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Request post successfully",
		"data": gin.H{
			"postData": post,
		},
	})
}
