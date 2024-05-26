package models

type Comment struct {
	ID     int
	Text   string
	Time   string //最新的修改时间
	Ban    bool   //是否被封禁
	PostID int
	UserID int
	User   User
}
