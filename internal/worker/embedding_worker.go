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
	Attempts  int
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

func (w *Worker) loop(id int) {
	defer w.wg.Done()
	for job := range w.queue {
		w.process(id, job)
	}
}

func (w *Worker) process(workerID int, job Job) {
	for {
		start := time.Now()
		err := w.handle(job)
		if err == nil {
			w.logger.Debug("job processed",
				slog.Int("worker", workerID),
				slog.String("memory_id", job.MemoryID),
				slog.Int("attempts", job.Attempts+1),
				slog.Duration("latency", time.Since(start)),
			)
			return
		}

		job.Attempts++
		if job.Attempts > w.maxRetries {
			// dead letter: retries exhausted, drop with a loud log
			w.logger.Error("job failed permanently",
				slog.String("memory_id", job.MemoryID),
				slog.Int("attempts", job.Attempts),
				slog.String("error", err.Error()),
			)
			return
		}

		delay := w.baseDelay * (1 << (job.Attempts - 1)) // exponential backoff
		w.logger.Warn("job failed, retrying",
			slog.String("memory_id", job.MemoryID),
			slog.Int("attempt", job.Attempts),
			slog.Duration("retry_in", delay),
			slog.String("error", err.Error()),
		)
		time.Sleep(delay)
	}
}

func (w *Worker) Enqueue(job Job) {
	w.Queue <- job
}
