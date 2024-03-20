package goque

import (
	"errors"
	"sync"
	"sync/atomic"
)

type JobQueue struct {
	queue  chan func()
	active atomic.Bool
	maxGo  int
	curGo  atomic.Int32
	wg     sync.WaitGroup
}

func NewJobQueue(maxSize int, maxGoroutine int) (*JobQueue, error) {
	if maxSize <= 1 {
		return nil, errors.New("buffer size must be > 1")
	}
	if maxGoroutine < 1 {
		return nil, errors.New("maxGoroutine must be >= 1")
	}
	return &JobQueue{
		queue: make(chan func(), maxSize),
		maxGo: maxGoroutine,
	}, nil
}

func (p *JobQueue) Start() {
	p.active.Store(true)
	go func() {

		for {

			if !p.active.Load() && len(p.queue) == 0 {
				break
			}

			if p.curGo.Load() < int32(p.maxGo) && len(p.queue) != 0 {
				p.curGo.Add(1)
				go func() {
					job := <-p.queue
					job()
					p.curGo.Add(-1)
					p.wg.Done()
				}()
			}

		}
	}()
}

func (p *JobQueue) Stop() {
	p.active.Store(false)
	p.wg.Wait()
}

func (p *JobQueue) Add(f func()) error {
	if !p.active.Load() {
		return errors.New("can't add new job because queue is inactive")
	}
	p.wg.Add(1)
	p.queue <- f
	return nil
}
