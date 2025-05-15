package event

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/command"
	"kmem/internal/utils"
	"time"
)

/*
	handle(context.Context) Result

	setResultChannel(chan Result)
	getResultChannel() chan Result

	setTimeout(time.Duration)
*/

type userLoggedIn struct {
	username string
	timeout  time.Duration
	resultCh chan Result
}

func UserLoggedIn(username string, options ...eventOption) *userLoggedIn {
	ul := &userLoggedIn{
		username: username,
		resultCh: nil,
		timeout:  defaultTimeout,
	}

	for _, opt := range options {
		opt(ul)
	}

	return ul
}

func (ul *userLoggedIn) setResultChannel(rchan chan Result) {
	ul.resultCh = rchan
}

func (ul userLoggedIn) getResultChannel() chan Result {
	return ul.resultCh
}

func (ul *userLoggedIn) setTimeout(t time.Duration) {
	ul.timeout = t
}

func (ul *userLoggedIn) handle(ctx context.Context, pg *database.Postgres) Result {
	accessToken, err := utils.CreateJwt(utils.ACCESSTOKEN_MAX_AGE, ul.username)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to create access token: %v", err), nil)
	}

	refreshToken, err := utils.CreateJwt(utils.REFRESHTOKEN_MAX_AGE, ul.username)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to create refresh token: %v", err), nil)
	}

	if err := command.UpdateLastLogin(pg, ul.username); err != nil {
		return newEventResult(utils.FAIL, "login fail", nil)
	}

	return newEventResult(utils.SUCCESS, "login success", map[string]string{
		utils.ACCESS_TOKEN:  accessToken,
		utils.REFRESH_TOKEN: refreshToken,
	})
}
