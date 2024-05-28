package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
)

type ForbidParam struct {
	UserId int `form:"userId" binding:"required"`
	//修改的是那种权限
	BanOption *int `form:"banOption" binding:"required"`
	//true 给帖子添加某种权限 false 给帖子撤销某种权限
	IsAdd *bool `form:"isAdd" binding:"required"`
}

// ForbidUser 操作用户权限
// @Tags 管理员方法
// @Summary 修改用户权限
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param userId formData int true "1"
// @Param banOption formData int true "0"
// @Param isAdd formData bool true "false"
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/forbid-user [post]
func ForbidUser(c *gin.Context) {
	var forbidParam ForbidParam
	if err := c.ShouldBind(&forbidParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	var userBan int
	err := models.DB.Model(&models.User{ID: forbidParam.UserId}).Select("ban").Scan(&userBan).Error
	if err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	//如果即将被修改的用户是管理员，那么修改将会失败(换个角度来看，管理员一旦给另一个用户管理员权限，那么后续将无法修改新增管理员的权限[最后就看谁有修改物理数据库的权限了])
	if models.JudgePermit(models.Admin, userBan) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The operation failed because the user is also an administrator",
		})
		return
	}
	if *forbidParam.IsAdd {
		userBan = models.AddPermit(*forbidParam.BanOption, userBan)
	} else {
		userBan = models.SubPermit(*forbidParam.BanOption, userBan)
	}
	err = models.DB.Model(&models.User{ID: forbidParam.UserId}).Update("ban", userBan).Error
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
		"msg":  "Succeeded in modifying the post permission",
	})
}
