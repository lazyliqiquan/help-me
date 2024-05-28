package click

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

type BothPostId struct {
	SeekHelpId int `form:"seekHelpId" binding:"required"`
	LendHandId int `form:"lendHandId" binding:"required"`
}

// AdoptHelp 采纳帮助者的帖子(一经采纳，无法撤销)
// @Tags 用户方法
// @Summary 采纳帮助者的帖子
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param seekHelpId formData int true "1"
// @Param lendHandId formData int true "2"
// @Success 200 {string} json "{"code":"0"}"
// @Router /adopt-help [post]
func AdoptHelp(c *gin.Context) {
	var bothPostId BothPostId
	if err := c.ShouldBind(&bothPostId); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	post := &models.Post{}
	err := models.DB.Model(&models.Post{ID: bothPostId.LendHandId}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "reward", "message")
	}).Select("id").First(post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Lend hand id is nonentity",
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
	//fixme Scan 一定要是一个结构体吗？
	var reward int
	err = models.DB.Model(&models.Post{ID: bothPostId.SeekHelpId}).Select("reward").Scan(&reward).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Seek help id is nonentity",
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
	post.User.Reward += reward
	post.User.Message = append(post.User.Message, post.ID)
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(models.Post{}).Where("id IN ?", []int{bothPostId.SeekHelpId, bothPostId.LendHandId}).
			Update("status", true).Error
		if err != nil {
			return err
		}
		//fixme Updates 里面的User，有必要加 & 吗
		return tx.Model(models.User{ID: post.UserID}).Select("reward", "message").Updates(&post.User).Error
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
		"msg":  "You have adopted this help post",
	})
}
