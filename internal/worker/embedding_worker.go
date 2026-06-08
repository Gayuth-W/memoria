package worker

type Job struct {
	MemoryID  string
	UserID    string
	SessionID string
	Text      string
}

type Worker struct {
	Queue  chan Job
	Handle func(Job) error
}

func NewWorker(buffer int, handler func(Job) error) *Worker {
	w := &Worker{
		Queue:  make(chan Job, buffer),
		Handle: handler,
	}

	go w.loop()
	return w
}

func (w *Worker) loop() {
	for job := range w.Queue {
		_ = w.Handle(job)
	}
}

func (w *Worker) Enqueue(job Job) {
	w.Queue <- job
}
