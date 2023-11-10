package libType

type Message struct {
	Id           int64
	ChatId       int64
	SenderUserId int64
	Text         string
	ContentType  string
}
