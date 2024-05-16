package models

type User struct {
	ID              int
	Name            string
	Avatar          string
	Reward          int
	Ban             int
	Email           string
	Password        string
	RegisterTime    string
	CommentSurplus  int    //评论剩余次数
	LatePublishDate string //最近一次发布的时间
	Message         GormIntList
	Collect         GormIntList
	Private         []Post
}

const (
	Admin int = iota
	Login
	PublishSeekHelp //发表求助
	ModifySeekHelp  //编辑求助
	PublishLendHand //发表帮助
	ModifyLendHand  //编辑帮助
	PublishComment
	ModifyComment
)

func JudgePermit(option, ban int) bool {
	return (ban & (1 << option)) == 0
}
