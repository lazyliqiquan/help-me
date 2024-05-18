package before

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Common
// 这里负责处理添加帖子和修改帖子的公共部分
func Common(isAdd bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//"0" 表示求助 "1"表示帮助
		postType := c.PostForm("postType")
		title := c.PostForm("title")
		createTime := c.PostForm("createTime")
		updateTime := c.PostForm("updateTime")
		language := c.PostForm("language")
		document := c.PostForm("document")
		imageNumStr := c.PostForm("imageNum")
		tags := c.PostForm("tags")
		if utils.IsNuiStrs(title, language, document, imageNumStr) || (isAdd && (createTime == "" || postType == "")) || (!isAdd && updateTime == "") {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Missing parameter",
			})
			c.Abort()
		}
		imageNum, err := strconv.Atoi(imageNumStr)
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The imageNum parameter is not a legal integer",
			})
			c.Abort()
		}
		//fixme 标签不能包含'|',由前端负责检测
		tagList := strings.Split(tags, "|")
		var imageFilesPath []string
		for i := 0; i < imageNum; i++ {
			file, _, err := c.Request.FormFile("image" + strconv.Itoa(i))
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Image file parsing error",
				})
				c.Abort()
			}
			imageFilesPath = append(imageFilesPath, config.Config.ImageFilePath+utils.GetUUID())
			// 反正存的是二进制，具体的文件类型问题应该不大吧
			err = utils.SaveAFile(imageFilesPath[i], file)
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Failed to save image file",
				})
				c.Abort()
			}
			//就算调用c.Abort()，也还是会执行defer的吧
			defer func(file string) {
				win := c.GetBool("win")
				if !win {
					if err := os.Remove(file); err != nil {
						utils.Logger.Errorln(err)
					}
				}
			}(imageFilesPath[i])
		}
		codeFile, _, err := c.Request.FormFile("codeFile")
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Code file parsing error",
			})
			c.Abort()
		}
		codeFilePath := config.Config.CodeFilePath + utils.GetUUID() + ".txt"
		err = utils.SaveAFile(codeFilePath, codeFile)
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Failed to save code file",
			})
			c.Abort()
		}
		defer func(file string) {
			//不存在就返回默认值，即false
			win := c.GetBool("win")
			if !win {
				if err := os.Remove(file); err != nil {
					utils.Logger.Errorln(err)
				}
			}
		}(codeFilePath)
		newPost := &models.Post{
			Title:      title,
			CreateTime: createTime,
			Language:   language,
			Tags:       tagList,
			PostStats: models.PostStats{
				CodePath:   codeFilePath,
				Document:   document,
				UpdateTime: updateTime,
				ImagePath:  imageFilesPath,
			},
		}
		if isAdd {
			c.Set("postType", postType)
			if postType == "0" {
				selectReward := c.GetInt("selectReward")
				newPost.Reward = selectReward
			}
		}
		c.Set("newPost", newPost)
		//调用c.Next()，还是会回到这里的，所以调用c.Next()的时候，不会执行defer，也就是不会删除文件
		c.Next()
	}
}
