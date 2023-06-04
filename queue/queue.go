package queue

import (
	"log"
	"sync"
)

type TaskQueue struct {
	mu sync.Mutex
	tasks []*Task
	users map[int64]*Task
	Limit int
	Count int
}

func NewTaskQueue(limit int) *TaskQueue {
	return &TaskQueue{
		tasks: make([]*Task, 0),
		users: make(map[int64]*Task, 0),
		Limit: limit,
	}
}


// Get task by UserID and its count in queue
func (q *TaskQueue) Load(userId int64) (*Task, int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for n, task := range q.tasks {
		if task.UserID == userId {
			return task, n
		}
	}

	return nil, -1
}


func (q *TaskQueue) Enqueue(task *Task) (int, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	t, exists := q.users[task.UserID]
	if exists {
		// update existing
		if t.MessageID == task.MessageID {
			t.Question = task.Question
			return q.Count, nil
		}

		return q.Count, ErrOnePerUser
	}

	if q.Count == q.Limit {
		return q.Count, ErrQueueLimit
	}


	q.tasks = append(q.tasks, task)
	q.users[task.UserID] = task
	q.Count++

	return q.Count, nil
}

func (q *TaskQueue) Dequeue() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.Count == 0 {
		return nil, ErrQueueEmpty
	}

	task := q.tasks[0]
	log.Println("task:", task)
	log.Println("q.Count:", q.Count)

	q.tasks[0] = nil
	q.tasks = q.tasks[1:]
	delete(q.users, task.UserID)

	q.Count--

	return task, nil
}