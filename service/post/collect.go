package post

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

type CollectParam struct {
	PostId int  `form:"postId" binding:"required"`
	IsAdd  bool `form:"isAdd" binding:"required"`
}

// UpdateCollect 修改用户收藏夹
// @Tags 用户方法
// @Summary 修改用户收藏夹
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param postId formData int true 1
// @Param isAdd formData bool true false
// @Success 200 {string} json "{"code":"0"}"
// @Router /update-collect [post]
func UpdateCollect(c *gin.Context) {
	var collectParam CollectParam
	if err := c.ShouldBind(&collectParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	userId := c.GetInt("id")
	collectList := make([]int, 0)
	err := models.DB.Model(&models.User{ID: userId}).Select("collect").Scan(&collectList).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	if collectParam.IsAdd {
		collectList = append(collectList, collectParam.PostId)
	} else {
		index := -1
		for i, v := range collectList {
			if v == collectParam.PostId {
				index = i
				break
			}
		}
		if index != -1 {
			collectList = append(collectList[0:index], collectList[index+1:]...)
		}
	}
	err = models.DB.Model(&models.User{ID: userId}).Update("collect", collectList).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
}
