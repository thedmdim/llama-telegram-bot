package queue

import "fmt"

type Result struct {
	Text string
	Err error
}

type Task struct {
	UserId         int64
	MessageId      int
	AnnounceId     int
	Question       string
	Stop        chan bool

	Stream chan string
	Result chan Result
}

func (t *Task) WrapInRoles() string {
	return "### User: answer my next question. " + t.Question + "\n### Assistant:"
}