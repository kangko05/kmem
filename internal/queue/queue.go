package queue

import (
	"context"
	"sync"
)

type Queue struct {
	ctx     context.Context
	list    chan item
	workers int
}

func New(ctx context.Context) *Queue {
	q := &Queue{
		ctx:     ctx,
		list:    make(chan item),
		workers: 16, // default
	}

	go q.run()

	return q
}

func (q *Queue) run() {
	var wg sync.WaitGroup
	wg.Add(q.workers)

	for range q.workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-q.ctx.Done():
					return
				case it := <-q.list:
					it.process()
					// handle list items
				}
			}
		}()
	}

	wg.Wait()
}

func (q *Queue) Add(item item) {
	q.list <- item
}
