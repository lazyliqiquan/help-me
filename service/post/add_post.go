package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// AddPost
// 前面的中间件已经帮我们判断了该帖子到底能不能创建，这里我们只需要新建就好,不需要判断
// 新建帖子，把求助和帮助合起来了
// 共有参数：帖子类型、标题、创建时间、语言、文档文本、文档图片集、图片数量、代码文件、标签
// 求助特有：悬赏、
// 帮助特有：求助id、
func AddPost(c *gin.Context) {
	userId := c.GetInt("id")
	seekHelpId := c.GetInt("seekHelpId")
	selectReward := c.GetInt("selectReward")
	reward := c.GetInt("reward")
	postType := c.PostForm("postType")
	title := c.PostForm("title")
	createTime := c.PostForm("createTime")
	language := c.PostForm("language")
	document := c.PostForm("document")
	imageNumStr := c.PostForm("imageNum")
	tags := c.PostForm("tags")
	if utils.IsNuiStrs(postType, title, createTime, language, document, imageNumStr) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Missing parameter",
		})
		return
	}
	imageNum, err := strconv.Atoi(imageNumStr)
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The imageNum parameter is not a legal integer",
		})
		return
	}
	//fixme 标签不能包含'|',由前端负责检测
	tagList := strings.Split(tags, "|")
	var imageFiles []multipart.File
	var imageFilesPath []string
	for i := 0; i < imageNum; i++ {
		file, _, err := c.Request.FormFile("image" + strconv.Itoa(i))
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Image file parsing error",
			})
			return
		}
		// 这里关闭的文件应该不总是最后一个吧(如果程序内存溢出可以考虑文件是否及时关闭)
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				utils.Logger.Errorln(err)
			}
		}(file)
		imageFiles = append(imageFiles, file)
		// 反正存的是二进制，具体的文件类型问题应该不大吧
		imageFilesPath = append(imageFilesPath, config.Config.ImageFilePath+utils.GetUUID())
	}
	codeFile, _, err := c.Request.FormFile("codeFile")
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Code file parsing error",
		})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			utils.Logger.Errorln(err)
		}
	}(codeFile)
	codeFilePath := config.Config.CodeFilePath + utils.GetUUID() + ".txt"
	newPost := &models.Post{
		Title:      title,
		CreateTime: createTime,
		Reward:     selectReward,
		Language:   language,
		Tags:       tagList,
		PostStats: models.PostStats{
			CodePath:  codeFilePath,
			Document:  document,
			ImagePath: imageFilesPath,
		},
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
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
		if err != nil {
			return err
		}
		//保存文件也放在事务里面，防止保存文件的时候报错，导致数据库和文件不一致
		for i, v := range imageFilesPath {
			if err := utils.SaveAFile(v, imageFiles[i]); err != nil {
				return err
			}
		}
		return utils.SaveAFile(codeFilePath, codeFile)
	})
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error or save file fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Add post successfully",
	})
}
