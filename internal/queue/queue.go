package queue

import (
	"context"
	"sync"
)

type QueueItemType uint

const (
	GEN_THUMBNAILS QueueItemType = iota
)

type QueueItem struct {
	itemType QueueItemType
}

func (it *QueueItem) process() {
	switch it.itemType {
	case GEN_THUMBNAILS:
		// gen thumbnails
	}
}

type Queue struct {
	ctx     context.Context
	list    chan QueueItem
	workers int
}

func New(ctx context.Context) *Queue {
	q := &Queue{
		ctx:     ctx,
		list:    make(chan QueueItem),
		workers: 16, // default
	}

	go q.run()

	return q
}

func (q *Queue) run() {
	var wg sync.WaitGroup

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

func (q *Queue) Add(itemType QueueItemType) {
	q.list <- QueueItem{
		itemType: itemType,
	}
}
