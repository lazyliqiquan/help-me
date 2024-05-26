package models

type Comment struct {
	ID     int
	Text   string
	Time   string //最新的修改时间
	PostID int
	UserID int
	User   User
}
