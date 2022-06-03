package model

type Dialog struct {
	Id            int
	Type          string
	Title         string
	LastMessage   string
	LastMessageId int64
	Platform      string
	Date          int
}

type Messages struct {
	DialogId     int
	MessagesList []Message
}

type Message struct {
	MessageId       int
	MessageText     string
	MessageIdAuthor int
	MessageAuthor   string
	MessageDate     int
}
