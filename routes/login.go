package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/middlewares/post/before"
	"github.com/lazyliqiquan/help-me/service"
	"github.com/lazyliqiquan/help-me/service/comment"
	"github.com/lazyliqiquan/help-me/service/post"
	"github.com/lazyliqiquan/help-me/service/user"
)

// 需要完成登录才可访问
func login(r *gin.Engine) {
	//在安全模式下，仅管理员可操作
	safeAuth := r.Group("/", middlewares.TokenSafeModel())
	{
		safeAuth.POST("/download-file", service.DownloadFile)
		safeAuth.POST("/logout-post-list", post.LogoutPostList)
	}
	viewAuth := safeAuth.Group("/view")
	{
		viewAuth.POST("/seek-help", middlewares.View(middlewares.SeekHelpItem), post.ViewPost)
		viewAuth.POST("/lend-hand", middlewares.View(middlewares.LendHandItem), post.ViewPost)
		viewAuth.POST("/comment", middlewares.View(middlewares.CommentItem), comment.View)
	}
	loginAuth := safeAuth.Group("/", middlewares.LoginModel())
	{
		loginAuth.POST("/private-post-list", post.PrivatePostList)
		loginAuth.POST("/collect-post-list", post.CollectPostList)
		loginAuth.POST("/update-collect", post.UpdateCollect)
	}
	modifyAuth := loginAuth.Group("/modify")
	{
		modifyAuth.POST("/seek-help", middlewares.Modify(middlewares.SeekHelpItem), before.ModifySeekHelp(), before.Common(false), post.ModifyPost)
		modifyAuth.POST("/lend-hand", middlewares.Modify(middlewares.LendHandItem), before.ModifyLendHand(), before.Common(false), post.ModifyPost)
		modifyAuth.POST("/comment", middlewares.Modify(middlewares.CommentItem), comment.Modify)
	}
	publishAuth := loginAuth.Group("/publish")
	{
		publishAuth.POST("/seek-help", middlewares.Publish(middlewares.SeekHelpItem), before.AddSeekHelp(), before.Common(true), post.AddPost)
		publishAuth.POST("/lend-hand", middlewares.Publish(middlewares.LendHandItem), before.AddLendHand(), before.Common(true), post.AddPost)
		publishAuth.POST("/comment", middlewares.Publish(middlewares.CommentItem), comment.Add)
	}
	adminAuth := loginAuth.Group("/admin", middlewares.AdminModel())
	{
		adminAuth.POST("/forbid-one-comment", comment.ForbidOneComment)
		adminAuth.POST("/forbid-post", post.ForbidPost)
		adminAuth.POST("/forbid-user", user.ForbidUser)
	}
}
