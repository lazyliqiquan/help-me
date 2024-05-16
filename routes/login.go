package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/middlewares/post/before"
	"github.com/lazyliqiquan/help-me/service"
	"github.com/lazyliqiquan/help-me/service/post"
)

// 需要完成登录才可访问
func login(r *gin.Engine) {
	//在安全模式下，仅管理员可操作
	safeAuth := r.Group("/", middlewares.TokenSafeModel())
	{
		safeAuth.POST("/download-file", service.DownloadFile)
	}
	viewAuth := safeAuth.Group("/view")
	{
		viewAuth.POST("/seek-help", middlewares.View(middlewares.SeekHelpItem))
		viewAuth.POST("/lend-hand", middlewares.View(middlewares.LendHandItem))
		viewAuth.POST("/comment", middlewares.View(middlewares.CommentItem))
	}
	loginAuth := safeAuth.Group("/", middlewares.LoginModel())
	modifyAuth := loginAuth.Group("/modify")
	{
		modifyAuth.POST("/seek-help", middlewares.Modify(middlewares.SeekHelpItem), before.ModifySeekHelp())
		modifyAuth.POST("/lend-hand", middlewares.Modify(middlewares.LendHandItem), before.ModifyLendHand())
		modifyAuth.POST("/comment", middlewares.Modify(middlewares.CommentItem))
	}
	publishAuth := loginAuth.Group("/publish")
	{
		publishAuth.POST("/seek-help", middlewares.Publish(middlewares.SeekHelpItem), before.AddSeekHelp(), post.AddPost)
		publishAuth.POST("/lend-hand", middlewares.Publish(middlewares.LendHandItem), before.AddLendHand(), post.AddPost)
		publishAuth.POST("/comment", middlewares.Publish(middlewares.CommentItem))
	}

}
