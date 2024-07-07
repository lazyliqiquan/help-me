package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// SubmitPost
// 关于求助帖子和帮助帖子的添加和修改都通过该函数来完成
// fixme 帖子上传过来以后，并没有检查它的文档大小和图片大小，只是在前端检查了，这样可以轻松突破约束
func SubmitPost(c *gin.Context) {
	postType := c.GetInt("postType")
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
	var err error
	if postType == 0 || postType == 2 {
		userId := c.GetInt("id")
		var selectReward, reward int
		if postType == 0 {
			reward = c.GetInt("reward")
			selectReward, err = strconv.Atoi(c.PostForm("reward"))
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The reward parameter is not a legal integer",
				})
				return
			}
			if reward <= 0 || selectReward > reward || selectReward <= 0 {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Illegality reward",
				})
				return
			}
			newPost.Reward = selectReward
		}
		//必须是事务，保证操作的原子性
		err = models.DB.Transaction(func(tx *gorm.DB) error {
			//创建之后，是会返回新创建的数据的id的吧
			err := tx.Model(&models.Post{}).Create(newPost).Error
			if err != nil {
				return err
			}
			err = tx.Model(&models.User{ID: userId}).Association("Private").Append(&models.Post{ID: newPost.ID})
			if err != nil {
				return err
			}
			if postType == 0 {
				//减去对应的悬赏
				err = tx.Model(&models.User{ID: userId}).Update("reward", reward-selectReward).Error
			} else {
				seekHelpId := c.GetInt("seekHelpId")
				tempPost := &models.Post{}
				err = tx.Model(&models.Post{ID: seekHelpId}).Preload("User", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "message")
				}).First(tempPost).Error
				if err != nil {
					return err
				}
				tempPost.User.Message = append(tempPost.User.Message, newPost.ID)
				err = tx.Model(&models.User{ID: tempPost.User.ID}).Update("message", tempPost.User.Message).Error
				if err != nil {
					return err
				}
				//将帮助帖子添加到对应的求助帖子的帮助列表下
				var lendHandSum int
				err = tx.Model(&models.Post{ID: seekHelpId}).Select("lend_hand_sum").Scan(&lendHandSum).Error
				if err != nil {
					return err
				}
				err = tx.Model(&models.Post{ID: seekHelpId}).Update("lend_hand_sum", lendHandSum+1).Error
				if err != nil {
					return err
				}
				err = tx.Model(&models.Post{ID: seekHelpId}).Association("LendHands").Append(&models.Post{ID: newPost.ID})
			}
			return err
		})
	} else {
		postId := c.GetInt("seekHelpId")
		if postType == 3 {
			postId = c.GetInt("lendHandId")
		}
		err = models.DB.Model(&models.Post{ID: postId}).Updates(newPost).Error
	}
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
		"msg":  "Edit post successfully",
	})
	//设置一个标志位，表示操作成功；操作成功不需要删除之前在中间件中创建的文件，操作失败则需要删除
	c.Set("win", true)
}
