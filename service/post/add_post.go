package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

// AddPost
// 前面的中间件已经帮我们判断了该帖子到底能不能创建，这里我们只需要新建就好,不需要判断
// 新建帖子，把求助和帮助合起来了
// 共有参数：帖子类型、标题、创建时间、语言、文档文本、文档图片集、图片数量、代码文件、标签
// 求助特有：悬赏、
// 帮助特有：求助id、
func AddPost(c *gin.Context) {
	userId := c.GetInt("id")
	postType := c.GetString("postType")
	seekHelpId := c.GetInt("seekHelpId")
	selectReward := c.GetInt("selectReward")
	reward := c.GetInt("reward")
	_newPost, _ := c.Get("newPost")
	newPost, ok := _newPost.(*models.Post)
	if !ok {
		utils.Logger.Errorln("assertion fail")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Assertion fail",
		})
		return
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		//创建之后，是会返回新创建的数据的id的吧
		err := tx.Model(&models.Post{}).Create(newPost).Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.User{ID: userId}).Association("Private").Append(&models.Post{ID: newPost.ID})
		if err != nil {
			return err
		}
		if postType == "0" {
			//减去对应的悬赏
			err = tx.Model(&models.User{ID: userId}).Update("reward", reward-selectReward).Error
		} else {
			//将帮助帖子添加到对应的求助帖子的帮助列表下
			err = tx.Model(&models.Post{ID: seekHelpId}).Association("LendHands").Append(&models.Post{ID: newPost.ID})
		}
		return err
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
		"msg":  "Add post successfully",
	})
	//设置一个标志位，表示操作成功；操作成功不需要删除之前在中间件中创建的文件，操作失败则需要删除
	c.Set("win", true)
}
