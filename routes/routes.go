package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help-me/middlewares"
	"github.com/lazyliqiquan/help-me/middlewares/post/before"
	"github.com/lazyliqiquan/help-me/service/click"
	"github.com/lazyliqiquan/help-me/service/comment"
	"github.com/lazyliqiquan/help-me/service/post"
	"github.com/lazyliqiquan/help-me/service/setting"
	"github.com/lazyliqiquan/help-me/service/user"
)

// 需要完成登录才可访问
func routes(r *gin.Engine) {
	//登录模块还是要单独分出来，免得管理员的token过期了，然后网站还处于安全模式，那就死锁了
	r.POST("/login", user.Login)
	//在安全模式下，仅管理员可操作
	safeAuth := r.Group("/", middlewares.TokenSafeModel())
	{
		safeAuth.POST("/send-code", user.SendCode)
		safeAuth.POST("/register", user.Register)
		safeAuth.POST("/find-password", user.FindPassword)
		safeAuth.Static("/files", "./files")
		//safeAuth.POST("/download-file", other.DownloadFile)
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
		loginAuth.POST("/adopt-help", click.AdoptHelp)
		loginAuth.POST("/upvote", click.Upvote)
		loginAuth.POST("/mark-single-info", click.MarkSingleInfo)
		loginAuth.POST("/mark-all-info", click.MarkAllInfo)
		loginAuth.POST("/message-list", post.MessageList)
		loginAuth.POST("/before-edit", post.BeforeEdit)
	}
	modifyAuth := loginAuth.Group("/modify")
	{
		modifyAuth.POST("/seek-help", middlewares.Modify(middlewares.SeekHelpItem), before.ModifySeekHelp(), before.Common(false), post.ModifyPost)
		modifyAuth.POST("/lend-hand", middlewares.Modify(middlewares.LendHandItem), before.ModifyLendHand(), before.Common(false), post.ModifyPost)
		modifyAuth.POST("/comment", middlewares.Modify(middlewares.CommentItem), comment.Modify)
	}
	publishAuth := loginAuth.Group("/publish")
	{
		publishAuth.POST("/seek-help", middlewares.Publish(middlewares.SeekHelpItem), before.Common(true), post.AddPost)
		publishAuth.POST("/lend-hand", middlewares.Publish(middlewares.LendHandItem), before.AddLendHand(), before.Common(true), post.AddPost)
		publishAuth.POST("/comment", middlewares.Publish(middlewares.CommentItem), comment.Add)
	}
	adminAuth := loginAuth.Group("/admin", middlewares.AdminModel())
	{
		adminAuth.POST("/forbid-one-comment", comment.ForbidOneComment)
		adminAuth.POST("/forbid-post", post.ForbidPost)
		adminAuth.POST("/forbid-user", user.ForbidUser)
		adminAuth.POST("/modify-permission", setting.ModifyPermission)
		adminAuth.POST("/view-permission", setting.ViewPermission)
		adminAuth.POST("/modify-restriction", setting.ModifyRestriction)
		adminAuth.POST("/view-restriction", setting.ViewRestriction)
	}
}
