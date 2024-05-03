package models

type Post struct {
	ID          int
	Title       string
	CreateTime  string
	Reward      int
	Language    string
	LikeSum     int
	LendHandSum int //可以通过该成员来判断是seekHelp(>=0)还是lendHand(-1)
	CommentSum  int
	Status      int
	Tags        GormStrList
	Ban         int
	UserID      int
	User        User
	Comments    []Comment
	PostID      int
	Posts       []Post `gorm:"foreignKey:postID"`
	PostStats   PostStats
}

type PostStats struct {
	ID         int
	UpdateTime string
	CodePath   string
	PageView   int
	Document   string
	ImagePath  GormStrList
	PostID     int
}
