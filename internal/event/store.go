package event

import (
	"context"
	"kmem/internal/database"
	"log"
	"sync"
)

type Store struct {
	ctx     context.Context
	evchan  chan eventType
	pg      *database.Postgres
	workers int
}

func NewStore(ctx context.Context, pg *database.Postgres) *Store {
	return &Store{
		ctx:     ctx,
		evchan:  make(chan eventType, 100),
		pg:      pg,
		workers: 8, // TODO: default 8 for now
	}
}

func (s *Store) Run() {
	log.Println("starting event store...")

	var wg sync.WaitGroup
	wg.Add(s.workers)

	for i := range s.workers {
		go s.runWorker(i, &wg)
	}

	<-s.ctx.Done()
	log.Println("context canceld, waiting for workers to stop...")

	wg.Wait()
	log.Println("all workers stopped, closing event store...")
}

func (s *Store) runWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("stopping worker %d", id)
			return

		case ev := <-s.evchan:
			res := ev.handle(s.ctx, s.pg)

			if resChan := ev.getResultChannel(); resChan != nil {
				resChan <- res
			}

			// TODO: store event result into db
		}
	}
}

func (s *Store) Register(evtype eventType, evOpts ...eventOption) {
	s.evchan <- evtype
}
