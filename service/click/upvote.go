package click

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

type LikeParam struct {
	PostId int   `form:"postId" binding:"required"`
	IsAdd  *bool `form:"isAdd" binding:"required"`
}

// Upvote 给某个帖子(求助或者帮助)点赞
// @Tags 用户方法
// @Summary 给某个帖子(求助或者帮助)点赞
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postId formData int true "1"
// @Param isAdd formData bool true "false"
// @Success 200 {string} json "{"code":"0"}"
// @Router /upvote [post]
func Upvote(c *gin.Context) {
	var likeParam LikeParam
	if err := c.ShouldBind(&likeParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		type Result struct {
			LikeSum int
			Likes   models.GormIntList
		}
		result := &Result{}
		err := tx.Model(&models.Post{ID: likeParam.PostId}).Select("like_sum", "likes").Scan(result).Error
		if err != nil {
			return err
		}
		userId := c.GetInt("id")
		if *likeParam.IsAdd {
			result.Likes = append(result.Likes, userId)
			result.LikeSum++
		} else {
			index := -1
			for i, v := range result.Likes {
				if v == userId {
					index = i
					break
				}
			}
			if index != -1 {
				result.Likes = append(result.Likes[0:index], result.Likes[index+1:]...)
				result.LikeSum--
			}
		}
		return tx.Model(&models.Post{ID: likeParam.PostId}).Select("like_sum", "likes").
			Updates(&models.Post{
				LikeSum: result.LikeSum,
				Likes:   result.Likes,
			}).Error
	})
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	msg := "Unlike successful"
	if *likeParam.IsAdd {
		msg = "Upvote successful"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
	})
}
