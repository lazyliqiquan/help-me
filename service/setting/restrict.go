package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/config"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"net/http"
	"time"
)

type RestrictionParam struct {
	Option   string `form:"option" binding:"required"`
	MaxLimit int    `form:"maxLimit" binding:"required"`
}

// ModifyRestriction 修改网站限制配置
// @Tags 管理员方法
// @Summary 修改网站限制配置
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Param option formData string true "maxDocumentHeight"
// @Param maxLimit formData int true 6000
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/modify-restriction [post]
func ModifyRestriction(c *gin.Context) {
	var restrictionParam RestrictionParam
	if err := c.ShouldBind(&restrictionParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	if err := models.RDB.Set(c, restrictionParam.Option, restrictionParam.MaxLimit, time.Duration(0)).Err(); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Redis error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Modifying the website restriction successfully.",
	})
}

// ViewRestriction 浏览网站限制配置
// @Tags 管理员方法
// @Summary 浏览网站限制配置
// @Accept multipart/form-data
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} json "{"code":"0"}"
// @Router /admin/view-restriction [post]
func ViewRestriction(c *gin.Context) {
	type Restriction struct {
		Name     string
		MaxLimit int
	}
	permissionList := make([]Restriction, 0)
	for k := range config.Config.GetRestrictionSetting() {
		v, err := models.RDB.Get(c, k).Int()
		if err != nil {
			utils.Logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis error",
			})
			return
		}
		permissionList = append(permissionList, Restriction{
			Name:     k,
			MaxLimit: v,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Obtaining the website restriction list succeeded",
		"data": gin.H{"list": permissionList},
	})
}
