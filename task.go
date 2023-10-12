package main

import "fmt"

type Task struct {
	UserID         int64
	MessageID      int
	AnnounceID     int
	Question       string
	Stopped        bool
	Stop        chan bool
}

func (t *Task) WrapInRoles(question string) {
	t.Question = fmt.Sprintf(SingleMessagePrompt, question)
}

func (t *Task) WrapPrevContext(previous, question string) {
	t.Question = fmt.Sprintf(ReplyMessagePrompt, previous, question)
}