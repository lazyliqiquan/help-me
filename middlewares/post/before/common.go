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
// 这里负责处理添加帖子和修改帖子的公共部分,下面是一些参数列表
// postType "0" 表示新建求助帖子 "1"表示新建帮助帖子 "2"表示修改帮助帖子 "3"表示修改求助帖子
// title 帖子的标题
// createTime 新建的帖子应该有该参数
// updateTime 修改的帖子应该有该参数
// document 文档的json形式
// tags 标签
// imageNameList 文档中新添加的图片的唯一标识符
// imageSizeList 文档中新添加的图片对应的大小
// originImageList 修改操作中，文档原本就有的图片，且没有删除
// seekHelpId 除了新建求助帖子都应该有
// lendHandId 只有修改帮助帖子有
// imageFiles 新添加的图片
func Common() gin.HandlerFunc {
	return func(c *gin.Context) {
		postType := c.PostForm("postType")
		title := c.PostForm("title")
		createTime := c.PostForm("createTime")
		updateTime := c.PostForm("updateTime")
		document := c.PostForm("document")
		tags := c.PostForm("tags")
		if utils.IsNuiStrs(postType, title, document) ||
			((postType == "0" || postType == "2") && createTime == "") ||
			((postType == "1" || postType == "3") && updateTime == "") {
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
		imageSizeList := strings.Split(c.PostForm("imageSizeList"), "|")
		if len(imageNameList) != len(imageSizeList) {
			utils.Logger.Errorln("The picture name list does not correspond to the picture size list")
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The picture name list does not correspond to the picture size list",
			})
			c.Abort()
		}
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
			}(imagePath)
			imagePath += "|" + imageSizeList[i]
			imageFilesPath = append(imageFilesPath, imagePath)
		}
		//修改帖子，主要是删除无用的图片
		if postType == "1" || postType == "3" {
			var postId int
			if postType == "1" {
				postId = c.GetInt("seekHelpId")
			} else {
				postId = c.GetInt("lendHandId")
			}
			oldPost := &models.Post{}
			err := models.DB.Model(&models.Post{ID: postId}).Preload("PostStats").First(oldPost).Error
			if err != nil {
				utils.Logger.Errorln(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Mysql error",
				})
				c.Abort()
			}
			//前端上传帖子资源之前，就应该把url图片资源解析出来
			originImageList := strings.Split(c.PostForm("originImageList"), "|")
			//收集无用的图片
			oldImageList := make([]string, 0)
			for _, v := range oldPost.PostStats.ImagePath {
				flag := false
				imagePath := strings.Split(v, "|")[0]
				for _, e := range originImageList {
					if imagePath == e {
						flag = true
						break
					}
				}
				if flag {
					//原本有用的图片继续留着
					imageFilesPath = append(imageFilesPath, v)
				} else {
					//没有用了的图片准备删除
					oldImageList = append(oldImageList, imagePath)
				}
			}
			defer func() {
				//如果后续操作失败，那么就不要删除图片
				if !c.GetBool("win") {
					return
				}
				for _, v := range oldImageList {
					if err := os.Remove(v); err != nil {
						//这里只需要打印到日志里面即可，不需要判定此次修改失败
						utils.Logger.Errorln(err)
					}
				}
			}()
		}
		//todo 后端也需要判断一下文档中的图片是否超出限制 imageFilesPath
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
