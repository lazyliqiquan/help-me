package post

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/models"
	"github.com/lazyliqiquan/help-me/utils"
	"gorm.io/gorm"
	"net/http"
)

// ListParam 因为请求的参数有点多，所以这里将请求的参数解析到准备好的结构体中
// fixme 结构体名称需要大写吗？
type ListParam struct {
	//请求类型:私人的求助列表、私人帮助列表、收藏求助列表、收藏帮助列表总共有六种情况
	ListType int `form:"listType"`
	//0 求助列表 >0 求助帖子对应的帮助列表
	SeekHelpId int `form:"seekHelpId"`
	//第几页，每页多少条
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"pageSize" binding:"required"`
	//过滤条件
	Status   int    `form:"status" binding:"required"`
	Language string `form:"language" binding:"required"`
	//排序条件，根据是否为true来判断使用那个条件来作为排序依据(用户一般只想看最新、最多点赞、最高悬赏、最多评论、最高活跃度的帖子)
	SortOption int `form:"sortOption" binding:"required"`
}

// LogoutPostList
// 帮助筛选条件：状态、语言。排序条件：日期、点赞、评论
// 求助筛选条件：状态、语言。排序条件：日期、点赞、评论、悬赏、活跃度
// @Tags 公共方法
// @Summary 请求帖子列表
// @Accept multipart/form-data
// @Param seekHelpId formData int true 0
// @Param page formData int true 1
// @Param pageSize formData int true 20
// @Param status formData int true 0
// @Param language formData string true "All"
// @Param sortOption formData int true 0
// @Success 200 {string} json "{"code":"0"}"
// @Router /view/login-post-list [post]
func LogoutPostList(c *gin.Context) {
	var listParam ListParam
	if err := c.ShouldBind(&listParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	tx := models.DB.Model(&models.Post{}).Session(&gorm.Session{})
	if listParam.SeekHelpId > 0 {
		seekHelpPost := &models.Post{}
		err := models.DB.Model(&models.Post{ID: listParam.SeekHelpId}).Preload("LendHands", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).First(&seekHelpPost).Error
		if errors.Is(err, gorm.ErrRecordNotFound) { //seekHelpId不存在
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Seek help id is nonentity",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql error",
			})
			return
		}
		lendHandList := make([]int, 0)
		for _, v := range seekHelpPost.LendHands {
			lendHandList = append(lendHandList, v.ID)
		}
		tx = tx.Where("id IN ? AND reward = ?", lendHandList, 0)
	} else {
		tx = tx.Where("reward > ?", 0)
	}
	//根据状态筛选
	if listParam.Status != 0 {
		tx = tx.Where("status = ?", listParam.Status == 1)
	}
	//根据语言筛选
	if listParam.Language != "All" {
		tx = tx.Where("language = ?", listParam.Language)
	}
	//先根据筛选条件得到总的数据量先
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	//排序
	sortCondition := "id"
	if listParam.SortOption == 1 {
		sortCondition = "like_sum"
	} else if listParam.SortOption == 2 {
		sortCondition = "comment_sum"
	} else if listParam.SortOption == 3 {
		sortCondition = "reward"
	} else if listParam.SortOption == 4 {
		sortCondition = "lend_hand_sum"
	}
	//先初始化，防止出现数组为nil的情况
	postList := make([]models.Post, 0)
	err := tx.Order(sortCondition+" DESC").Offset((listParam.Page-1)*listParam.PageSize).
		Limit(listParam.PageSize).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "avatar")
	}).Find(&postList).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql error",
		})
		return
	}
	//判断一下帖子是否可以浏览
	userId := c.GetInt("id")
	userBan := c.GetInt("ban")
	for i, v := range postList {
		if userId == 0 && models.JudgePermit(models.View) { //未登录

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Get logout post list successfully",
		"data": gin.H{"total": 0, "list": postList},
	})
}

// LoginPostList
// 帮助筛选条件：状态、语言。排序条件：日期、点赞、评论
// 求助筛选条件：状态、语言。排序条件：日期、点赞、评论、悬赏、活跃度
// @Tags 公共方法
// @Summary 请求帖子列表
// @Accept multipart/form-data
// @Param listType formData int true 0
// @Param page formData int true 1
// @Param pageSize formData int true 20
// @Param status formData int true 0
// @Param language formData string true "All"
// @Param sortOption formData int true 0
// @Success 200 {string} json "{"code":"0"}"
// @Router /view/login-post-list [post]
func LoginPostList(c *gin.Context) {
	var listParam ListParam
	if err := c.ShouldBind(&listParam); err != nil {
		utils.Logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parse request params fail",
		})
		return
	}
	userId := c.GetInt("id")
	tx := models.DB.Model(&models.Post{})
}
