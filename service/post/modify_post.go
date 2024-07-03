package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

// ModifyPost
// 修改帖子，把求助和帮助合起来了
// 共有参数：帖子类型、帖子id，标题、修改时间、语言、文档文本、文档图片集、图片数量、代码文件、标签
// 求助特有：
// 帮助特有：
func ModifyPost(c *gin.Context) {
	postType := c.GetString("postType")
	_newPost, _ := c.Get("newPost")
	postId := c.GetInt("seekHelpId")
	newPost, ok := _newPost.(*models.Post)
	if !ok {
		utils.Logger.Errorln("assertion fail")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Assertion fail",
		})
		return
	}
	if postType == "3" {
		postId = c.GetInt("lendHandId")
	}
	err := models.DB.Model(&models.Post{ID: postId}).Updates(newPost).Error
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
		"msg":  "Modify post successfully",
	})
	c.Set("win", true)
}
