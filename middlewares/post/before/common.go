package before

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"math/rand"
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
		document := c.PostForm("document")
		tags := c.PostForm("tags")
		if utils.IsNuiStrs(postType, title, document) || (isAdd && createTime == "") || (!isAdd && updateTime == "") {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Missing parameter",
			})
			c.Abort()
		}
		//fixme 标签不能包含'|',由前端负责检测
		tagList := strings.Split(tags, "|")
		//在前端的时候，先用唯一标识符给图片起名字，具体的图片路径到后端的时候再构造,通过'|'分隔
		imageNameList := strings.Split(c.PostForm("imageNameList"), "|")
		var imageFilesPath []string
		for i, v := range imageNameList {
			//需要给图片加上文件类型后缀吗，不加quill能分辨出来吗
			imagePath := config.Config.ImageFilePath + "/" +
				strconv.Itoa(rand.Intn(config.Config.SecondImageDirAmount)) + "/" + utils.GetUUID()
			document = strings.Replace(document, v, imagePath, 1)
			file, _, err := c.Request.FormFile("image" + strconv.Itoa(i))
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Image file parsing error",
				})
				c.Abort()
			}
			imageFilesPath = append(imageFilesPath, imagePath)
			err = utils.SaveAFile(imagePath, file)
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
		newPost := &models.Post{
			Title:      title,
			CreateTime: createTime,
			Tags:       tagList,
			PostStats: models.PostStats{
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
