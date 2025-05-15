package event

import (
	"context"
	"kmem/internal/database"
	"log"
)

type Store struct {
	ctx    context.Context
	evchan chan eventType
	pg     *database.Postgres
}

func NewStore(ctx context.Context, pg *database.Postgres) *Store {
	return &Store{
		ctx:    ctx,
		evchan: make(chan eventType, 100),
		pg:     pg,
	}
}

func (s *Store) Run() {
	log.Println("starting event store...")

	for {
		select {
		case <-s.ctx.Done():
			log.Println("closing event store...")
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
