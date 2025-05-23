package event

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/utils"
	"time"
)

// type eventType interface {
// 	handle(context.Context, *database.Postgres, *database.Cache) Result
//
// 	setResultChannel(chan Result)
// 	getResultChannel() chan Result
//
// 	setTimeout(time.Duration)
// }

type userLoggedOut struct {
	username string
	resultCh chan Result
	timeout  time.Duration
}

func UserLoggedOut(username string, options ...eventOption) *userLoggedOut {
	ulo := &userLoggedOut{
		resultCh: nil,
		timeout:  defaultTimeout,
	}

	for _, opt := range options {
		opt(ulo)
	}

	return ulo
}

func (lo *userLoggedOut) setResultChannel(rchan chan Result) {
	lo.resultCh = rchan
}

func (lo *userLoggedOut) getResultChannel() chan Result {
	return lo.resultCh
}

func (lo *userLoggedOut) setTimeout(t time.Duration) {
	lo.timeout = t
}

func (lo *userLoggedOut) handle(ctx context.Context, pg *database.Postgres, cache *database.Cache) Result {
	cache.InvalidateUserCache(lo.username)
	return newEventResult(utils.SUCCESS, fmt.Sprintf("%s logout success", lo.username), nil)
}
