package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// BeforeAdd
// 添加评论之前，先判断今日评论次数是否足够
func BeforeAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		postId, err := strconv.Atoi(c.PostForm("postId"))
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Post id is not a integer",
			})
			c.Abort()
		}
		isBefore := c.PostForm("isBefore")
		userId := c.GetInt("id")
		post := &models.Post{}
		err = models.DB.Model(&models.Post{ID: postId}).Select("id", "ban", "comment_sum").First(post).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Post id is nonentity",
			})
			c.Abort()
		} else if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			c.Abort()
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
				c.Abort()
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
				c.Abort()
			}
			curDate := utils.GetCurrentDate()
			if user.CommentSurplus <= 0 && user.LatePublishDate == curDate {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "You have reached your maximum number of comments posted today",
				})
				c.Abort()
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
					c.Abort()
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
				c.Abort()
			}
		}
		maxCommentWords, err := models.RDB.Get(c, "maxCommentWords").Int()
		if isBefore == "0" {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "You are eligible to comment",
				"data": gin.H{"maxCommentWords": maxCommentWords},
			})
			c.Abort()
		} else {
			text := c.PostForm("text")
			if len(text) > maxCommentWords {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Your comment is out of limit",
				})
				c.Abort()
			}
			c.Set("text", text)
		}
		utils.Logger.Infoln("Comment before add")
		c.Set("postId", post.ID)
		c.Set("commentSum", post.CommentSum)
		c.Next()
	}
}
