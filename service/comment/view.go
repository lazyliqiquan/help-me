package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

type RequestComment struct {
	PostId   int `form:"postId" binding:"required"`
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"pageSize" binding:"required"`
}

// View
// 1. 若网站不允许浏览评论，该用户是否是管理员
// 2. 若网站需要登录方可浏览评论，用户是否登录
// 3. 传递到后端的帖子id是否存在
// 4. 若帖子不允许浏览评论，该用户是否是管理员
// @Tags 公共方法
// @Summary 浏览帖子所对应的评论
// @Accept multipart/form-data
// @Param postId formData int true 1
// @Param page formData int true 1
// @Param pageSize formData int true 20
// @Success 200 {string} json "{"code":"0"}"
// @Router /view/comment [post]
func View(c *gin.Context) {
	var requestComment RequestComment
	if err := c.ShouldBind(&requestComment); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	var total int64
	post := &models.Post{}
	err := models.DB.Model(&models.Post{ID: requestComment.PostId}).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			//fixme 这里可以圈套吗
			return db.Count(&total).Preload("User", func(db2 *gorm.DB) *gorm.DB {
				return db2.Select("id", "avatar")
			}).Offset((requestComment.Page - 1) * requestComment.PageSize).Limit(requestComment.PageSize)
		}).Select("id").First(&post).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	//初始化
	if post.Comments == nil {
		post.Comments = make([]models.Comment, 0)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "The list of comments was obtained successfully",
		"data": gin.H{"total": total, "list": post.Comments},
	})
}
