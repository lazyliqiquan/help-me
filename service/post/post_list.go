package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

// ListParam 因为请求的参数有点多，所以这里将请求的参数解析到准备好的结构体中
// fixme 结构体名称需要大写吗？
type ListParam struct {
	//0表示求助列表 >0表示对应的帮助列表
	SeekHelpId int `json:"seekHelpId" form:"seekHelpId"`
	//过滤条件
	Status   int    `json:"status" form:"status"`
	Language string `json:"language" form:"language"`
	//排序条件，根据是否为true来判断使用那个条件来作为排序依据(用户一般只想看最新、最多点赞、最高悬赏、最多评论、最高活跃度的帖子)
	SortOption int `json:"sortOption" form:"sortOption"`
}

// GetPostList	请求帖子列表
// 求助筛选条件：状态、语言。排序条件：日期、点赞、悬赏、评论、活跃度
// 帮助筛选条件：状态、语言。排序条件：日期、点赞、悬赏、评论
// @Tags 公共方法
// @Summary 请求帖子列表
// @Accept multipart/form-data
// @Param seekHelpId formData string true "0"
// @Param status formData string true "all"
// @Param language formData string true "all"
// @Param sortOption formData string true "0"
// @Success 200 {string} json "{"code":"0"}"
// @Router /view/post-list [post]
func GetPostList(c *gin.Context) {
	var listParam ListParam
	if err := c.ShouldBindJSON(&listParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
}
