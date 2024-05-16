package post

import "github.com/gin-gonic/gin"

// ModifyPost
// 修改帖子，把求助和帮助合起来了
// 共有参数：帖子类型、标题、修改时间、语言、文档文本、文档图片集、图片数量、代码文件、标签
// 求助特有：
// 帮助特有：
func ModifyPost(c *gin.Context) {
	userId := c.GetInt("id")
	postType := c.PostForm("postType")
	title := c.PostForm("title")
	updateTime := c.PostForm("updateTime")
	language := c.PostForm("language")
	document := c.PostForm("document")
	imageNumStr := c.PostForm("imageNum")
	tags := c.PostForm("tags")
}
