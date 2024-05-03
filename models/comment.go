package models

type Comment struct {
	ID       int
	Text     int
	SendTime string
	PostID   int
	UserID   int
	User     User
}
