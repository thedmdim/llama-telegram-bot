package queue

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

func (t *Task) WrapInRoles(question string) {
	t.Question = "### User: Response to my next request. " + question + "\n### Assistant:"
}

func (t *Task) WrapPrevContext(previous, question string) {
	t.Question = "### Assistant: " + previous + "\n### User: " + question + "\n### Assistant:"
}