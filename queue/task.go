package queue

type Task struct {
	UserID         int64
	MessageID      int
	AnnounceID     int
	Question       string
	Stopped        bool
	Stop        chan bool
}

func (t *Task) WrapInRoles(question string) {
	t.Question = "### User: Response to my next request. " + question + "\n### Assistant:"
}

func (t *Task) WrapPrevContext(previous, question string) {
	t.Question = "### Assistant: " + previous + "\n### User: " + question + "\n### Assistant:"
}