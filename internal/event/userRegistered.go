package event

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/command"
	"kmem/internal/models"
	"kmem/internal/utils"
	"time"
)

type userRegistered struct {
	user       models.User
	resultChan chan Result
	timeout    time.Duration
}

func UserRegistered(user models.User, options ...eventOption) *userRegistered {
	ur := &userRegistered{user: user, resultChan: nil, timeout: defaultTimeout}

	for _, opt := range options {
		opt(ur)
	}

	return ur
}

func (ur *userRegistered) handle(ctx context.Context, pg *database.Postgres) Result {
	if err := command.InsertUser(pg, ur.user); err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to store user: %v", err), nil)
	}

	return newEventResult(utils.SUCCESS, "successfully stored user: "+ur.user.Username, nil)
}

func (ur *userRegistered) setResultChannel(resultChan chan Result) {
	ur.resultChan = resultChan
}

func (ur userRegistered) getResultChannel() chan Result {
	return ur.resultChan
}

func (ur *userRegistered) setTimeout(t time.Duration) {
	ur.timeout = t
}
