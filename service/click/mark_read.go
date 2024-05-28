package click

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// MarkSingleInfo 将一条消息标记为已读
// @Tags 用户方法
// @Summary 将一条消息标记为已读
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postId formData int true "1"
// @Success 200 {string} json "{"code":"0"}"
// @Router /mark-single-info [post]
func MarkSingleInfo(c *gin.Context) {
	postId, err := strconv.Atoi(c.PostForm("postId"))
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Post id is illegality",
		})
		return
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		userId := c.GetInt("id")
		user := &models.User{}
		err := tx.Model(&models.User{ID: userId}).Select("message").First(user).Error
		if err != nil {
			return err
		}
		index := -1
		for i, v := range user.Message {
			if v == postId {
				index = i
				break
			}
		}
		if index != -1 {
			user.Message = append(user.Message[0:index], user.Message[index+1:]...)
		}
		return tx.Model(&models.User{ID: userId}).Update("message", user.Message).Error
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
		"msg":  "Mark as read successful",
	})
}

// MarkAllInfo 将所有消息标记为已读
// @Tags 用户方法
// @Summary 将所以消息标记为已读
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} json "{"code":"0"}"
// @Router /mark-all-info [post]
func MarkAllInfo(c *gin.Context) {
	userId := c.GetInt("id")
	err := models.DB.Model(&models.User{ID: userId}).Update("message", models.GormIntList{}).Error
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
		"msg":  "All information is marked as read",
	})
}
