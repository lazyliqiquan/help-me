package models

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	Avatar       string
	Reward       int
	RegisterTime string
	Ban          int
	Message      GormIntList
	Posts        []Post
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
