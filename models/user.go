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
	Posts        []int
}

const (
	Admin int = iota
	Login
	PublishSeekHelp //发表求助
	EditSeekHelp    //编辑求助
	PublishLendHand //发表帮助
	EditLendHand    //编辑帮助
	PublishComment
)

func JudgePermit(option, ban int) bool {
	return (ban & (1 << option)) == 0
}
