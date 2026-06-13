package worker

import (
	"log/slog"
	"sync"
	"time"
)

type Job struct {
	MemoryID  string
	UserID    string
	SessionID string
	Text      string
}

type Worker struct {
	queue      chan Job
	handle     func(Job) error
	workers    int
	maxRetries int
	baseDelay  time.Duration
	logger     *slog.Logger
	wg         sync.WaitGroup
}
type Config struct {
	Buffer     int
	Workers    int
	MaxRetries int
	BaseDelay  time.Duration
	Logger     *slog.Logger
}

func NewWorker(cfg Config, handler func(Job) error) *Worker {
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	if cfg.BaseDelay <= 0 {
		cfg.BaseDelay = 200 * time.Millisecond
	}
	w := &Worker{
		queue:      make(chan Job, cfg.Buffer),
		handle:     handler,
		workers:    cfg.Workers,
		maxRetries: cfg.MaxRetries,
		baseDelay:  cfg.BaseDelay,
		logger:     cfg.Logger,
	}
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)  //Adding new workder
		go w.loop(i) //Running the added new worker
	}
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
