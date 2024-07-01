package models

type Post struct {
	ID          int
	Title       string
	CreateTime  string
	Reward      int //可以通过是否大于0来判断是求助帖子还是帮助帖子
	LikeSum     int
	LendHandSum int
	CommentSum  int
	Status      bool
	Tags        GormStrList
	Likes       GormIntList //存放的是点赞的用户id
	Ban         int
	UserID      int
	User        User
	Comments    []Comment
	LendHandID  int
	LendHands   []Post `gorm:"foreignKey:LendHandID"`
	PostStats   PostStats
}

type PostStats struct {
	ID         int
	UpdateTime string
	PageView   int
	Document   string
	ImagePath  GormStrList //删除或者修改帖子的时候，还要知道该帖子包含哪些图片，便于删除
	PostID     int
}

const (
	View int = iota
	Modify
	ViewComment
	AddComment
	AddLendHand
)
