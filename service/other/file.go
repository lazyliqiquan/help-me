package other

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"net/http"
	"strings"
)

// DownloadFile @Tags 公共方法
// @Tags 公共方法
// @Summary 下载指定路径的文件
// @Accept multipart/form-data
// @Param filePath formData string true "filePath"
// @Success 200 {string} json "{"code":"0"}"
// @Router /download-file [post]
func DownloadFile(c *gin.Context) {
	filePath := c.PostForm("filePath")
	if msg, legal := judgeLegalPath(filePath); !legal {
		c.JSON(http.StatusOK, gin.H{
			"code": 1, "msg": msg,
		})
		return
	}
	c.File(filePath)
}

// judgeLegalPath 判断传递过来的文件路径是否合法
func judgeLegalPath(filePath string) (string, bool) {
	if filePath == "" {
		return "The path to the file cannot be empty", false
	}
	list := strings.Split(filePath, "/")
	if len(list) != 3 {
		return "The file path is not level 3", false
	} else if list[0]+"/" != config.Config.RootFilePath {
		return "The root directory does not match", false
	}
	for _, v := range []string{
		config.Config.ImageFilePath,
		config.Config.CodeFilePath,
		config.Config.AvatarFilePath} {
		if list[1]+"/" == v {
			return "", true
		}
	}
	return "The secondary directory does not match", false
}
