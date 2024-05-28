package before

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"os"
	"os/exec"
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
		if utils.IsNuiStrs(postType, title, language, document, imageNumStr) || (isAdd && createTime == "") || (!isAdd && updateTime == "") {
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
		//如果是帮助帖子，还需要和求助帖子比较
		if postType == "1" {
			seekHelpPost := &models.Post{}
			seekHelpId := c.GetInt("seekHelpId")
			err := models.DB.Model(&models.Post{ID: seekHelpId}).Preload("PostStats").First(seekHelpPost).Error
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Mysql error",
				})
				return
			}
			//求助帖子和帮助帖子是同一种语言才有比较的意义
			if seekHelpPost.Language == language {
				pathList := strings.Split(codeFilePath, "/")
				pathList[len(pathList)-1] = "diff_" + pathList[len(pathList)-1]
				diffFilePath := strings.Join(pathList, "/")
				cmd := exec.Command("sh", "-c", "diff -U 9999 "+seekHelpPost.PostStats.CodePath+" "+codeFilePath+" > "+diffFilePath)
				// 这里是同步，有点耗时间，可以考虑Start和Wait的异步结合
				// 虽然报错exit status 1，但是结果还是好的，所以感觉不用管这里的报错
				cmd.Run()
				defer func(originFile, diffFile string) {
					//不存在就返回默认值，即false
					win := c.GetBool("win")
					if !win {
						if err := os.Remove(diffFile); err != nil {
							utils.Logger.Errorln(err)
						}
					} else {
						if err := os.Remove(originFile); err != nil {
							utils.Logger.Errorln(err)
						}
					}
				}(codeFilePath, diffFilePath)
				codeFilePath = diffFilePath
			}
		}
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
		c.Set("postType", postType)
		c.Set("newPost", newPost)
		//fixme 调用c.Next()，还是会回到这里的，所以调用c.Next()的时候，不会执行defer，也就是不会删除文件
		c.Next()
	}
}
