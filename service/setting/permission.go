package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"time"
)

type PermissionParam struct {
	//修改权限类型
	Option string `form:"option" binding:"required"`
	//true 表示允许 false 表示不允许
	Ban bool `form:"ban" binding:"required"`
}

// ModifyPermission 修改网站权限配置
// @Tags 管理员方法
// @Summary 修改网站权限配置
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param option formData string true "safeBan"
// @Param ban formData bool true false
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/modify-permission [post]
func ModifyPermission(c *gin.Context) {
	var permissionParam PermissionParam
	if err := c.ShouldBind(&permissionParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	v := utils.Permit
	if !permissionParam.Ban {
		v = utils.Forbid
	}
	if err := models.RDB.Set(c, permissionParam.Option, v, time.Duration(0)).Err(); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Redis error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Modifying the website permission successfully.",
	})
}

// ViewPermission 浏览网站权限配置
// @Tags 管理员方法
// @Summary 浏览网站权限配置
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/view-permission [post]
func ViewPermission(c *gin.Context) {
	type Permission struct {
		Name string
		Ban  bool
	}
	permissionList := make([]Permission, 0)
	for k := range config.Config.GetPermissionSetting() {
		v, err := models.RDB.Get(c, k).Result()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis error",
			})
			return
		}
		permissionList = append(permissionList, Permission{
			Name: k,
			Ban:  v == utils.Permit,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Obtaining the website permission list succeeded",
		"data": gin.H{"list": permissionList},
	})
}
