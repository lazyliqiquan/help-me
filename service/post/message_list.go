package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

// MessageList 消息列表
// @Tags 用户方法
// @Summary 消息列表
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} json "{"code":"0"}"
// @Router /message-list [post]
func MessageList(c *gin.Context) {
	userId := c.GetInt("id")
	user := &models.User{}
	postList := make([]models.Post, 0)
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.User{ID: userId}).Select("message").First(user).Error
		if err != nil {
			return err
		}
		return tx.Model(&models.Post{}).Where("id IN ?", user.Message).Select("reward").Find(&postList).Error
	})
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	type Result struct {
		PostId     int
		IsSeekHelp bool
	}
	resultList := make([]Result, 0)
	for i, v := range postList {
		result := Result{PostId: user.Message[i]}
		if v.Reward > 0 {
			result.IsSeekHelp = true
		}
		resultList = append(resultList, result)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Loading information successfully",
		"data": gin.H{"list": resultList},
	})
}
