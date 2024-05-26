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
	Admin           int = iota //管理员权限
	Login                      //连登录权限都没有，就表示封禁状态
	PublishSeekHelp            //发表求助
	ModifySeekHelp             //编辑求助
	PublishLendHand            //发表帮助
	ModifyLendHand             //编辑帮助
	PublishComment             //发布评论
	ModifyComment              //修改评论
)

// JudgePermit 判断是否具有某种权限
func JudgePermit(option, ban int) bool {
	return (ban & (1 << option)) == 0
}

// SubPermit 删除某种权限
func SubPermit(option, ban int) int { return ban | (1 << option) }

// AddPermit 添加某种权限
func AddPermit(option, ban int) int { return ban &^ (1 << option) }
