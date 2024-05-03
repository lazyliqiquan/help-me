package before

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/service"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ModifySeekHelp
// @Tags 用户方法
// @Summary 检查该用户是否具有修改求助帖子的资格
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Param seekHelpId formData string true "seekHelpId"
// @Success 200 {string} json "{"code":"0"}"
// @Router /before/add-seek-help [post]
func ModifySeekHelp(c *gin.Context) {
	userId := c.GetInt("id")
	post := &models.Post{}
	var err error
	post.ID, err = strconv.Atoi(c.PostForm("seekHelpId"))
	if err != nil {
		service.Logger.Errorln(err)
		// 不是整数的情况应该交给前端处理，我们不需要额外说明
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Seek help id nonentity",
		})
		return
	}
	err = models.DB.Model(&models.Post{}).Preload("User").Where("id = ?", post.ID).Select("ban", "lend_hand_sum").Find(post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //传递过来的seekHelpId不存在，用户输入的url有问题
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Url error : seek help id not exist",
			})
			return
		}
		service.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failed",
		})
		return
	}
	if !models.JudgePermit(models.Admin, userId) {
		if userId != post.User.ID {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not the owner of the post and cannot modify it",
			})
			return
		}
		if post.LendHandSum > 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The current seek help post has a sponsor and cannot be modified",
			})
			return
		}
		if !models.JudgePermit(models.Modify, post.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The current seek help post cannot be modified",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "You are ready to start editing",
	})
}
