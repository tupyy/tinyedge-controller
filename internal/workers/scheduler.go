package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Worker interface {
	Do(ctx context.Context) error
	Name() string
}

type Scheduler struct {
	heartbeat time.Duration
	workers   []Worker
}

func New(heartbeat time.Duration) *Scheduler {
	return &Scheduler{
		heartbeat: heartbeat,
		workers:   make([]Worker, 0),
	}
}

func (s *Scheduler) AddWorker(w Worker) *Scheduler {
	s.workers = append(s.workers, w)
	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	work := make(chan struct{}, 1)
	ticker := time.NewTicker(s.heartbeat)

	doWork := func(ch chan struct{}) {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	zap.S().Info("scheduler started")
	for {
		select {
		case <-work:
			for _, w := range s.workers {
				err := w.Do(ctx)
				if err != nil {
					zap.S().Errorw("worker finished with error", "name", w.Name(), "error", err)
				}
			}
		case <-ticker.C:
			doWork(work)
		case <-ctx.Done():
			zap.S().Info("closing scheduler")
			return
		}
	}
}
