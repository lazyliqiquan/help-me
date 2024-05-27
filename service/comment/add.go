package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

type NewComment struct {
	//添加评论 : 帖子id | 修改评论 : 评论id
	Id   int    `form:"id" binding:"required"`
	Text string `form:"text" binding:"required"`
	//发送时间或更新时间
	Time string `form:"time" binding:"required"`
}

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
// @Param postId formData int true 1
// @Param text formData string true "Hello"
// @Param sendTime formData string true "2024-05-26 15:10:00"
// @Success 200 {string} json "{"code":"0"}"
// @Router /publish/comment [post]
func Add(c *gin.Context) {
	var newComment NewComment
	if err := c.ShouldBind(&newComment); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	userId := c.GetInt("id")
	post := &models.Post{}
	err := models.DB.Model(&models.Post{ID: newComment.Id}).Select("id", "ban", "comment_sum").First(post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Post id is nonentity",
		})
		return
	} else if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	userBan := c.GetInt("ban")
	//管理员可以无限发评论
	if !models.JudgePermit(models.Admin, userBan) {
		//判断当前帖子是否可以评论
		if !models.JudgePermit(models.AddComment, post.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The current post cannot be commented",
			})
			return
		}
		//判断用户当日的评论是否达到上限
		user := &models.User{}
		err = models.DB.Model(&models.User{ID: userId}).Select("comment_surplus", "late_publish_date").First(user).Error
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			return
		}
		curDate := utils.GetCurrentDate()
		if user.CommentSurplus <= 0 && user.LatePublishDate == curDate {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You have reached your maximum number of comments posted today",
			})
			return
		}
		//重置每天的评论次数
		if user.LatePublishDate != curDate {
			commentSurplus, err := models.RDB.Get(c, "dayUserCommentLimit").Int()
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis error",
				})
				return
			}
			//可能后面添加评论失败，这里又提前减一了，导致数据与操作不一致
			user.CommentSurplus = commentSurplus - 1
			user.LatePublishDate = curDate
		} else { //更新每天的评论次数
			//可能后面添加评论失败，这里又提前减一了，导致数据与操作不一致
			user.CommentSurplus--
		}
		err = models.DB.Model(&models.User{ID: userId}).
			Select("comment_surplus", "late_publish_date").Updates(user).Error
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			return
		}
	}
	comment := &models.Comment{
		Time: newComment.Time,
		Text: newComment.Text,
		//这里直接添加，关联模式会自动维护吗
		PostID: newComment.Id,
		UserID: userId,
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.Comment{}).Create(comment).Error
		if err != nil {
			return err
		}
		return tx.Model(&models.Post{ID: newComment.Id}).Update("comment_sum", post.CommentSum+1).
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
