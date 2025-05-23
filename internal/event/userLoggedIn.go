package event

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/command"
	"kmem/internal/database/query"
	"kmem/internal/utils"
	"log"
	"time"
)

type userLoggedIn struct {
	username  string
	jwtSecret string
	timeout   time.Duration
	resultCh  chan Result
}

func UserLoggedIn(secretKey, username string, options ...eventOption) *userLoggedIn {
	ul := &userLoggedIn{
		username:  username,
		jwtSecret: secretKey,
		resultCh:  nil,
		timeout:   defaultTimeout,
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

func (ul *userLoggedIn) handle(ctx context.Context, pg *database.Postgres, cache *database.Cache) Result {
	accessToken, err := utils.CreateJwt(utils.ACCESSTOKEN_MAX_AGE, ul.username, ul.jwtSecret)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to create access token: %v", err), nil)
	}

	refreshToken, err := utils.CreateJwt(utils.REFRESHTOKEN_MAX_AGE, ul.username, ul.jwtSecret)
	if err != nil {
		return newEventResult(utils.FAIL, fmt.Sprintf("failed to create refresh token: %v", err), nil)
	}

	if err := command.UpdateLastLogin(pg, ul.username); err != nil {
		return newEventResult(utils.FAIL, "login fail", nil)
	}

	// try to load files into cache
	go func() {
		userfiles, err := query.QueryUserFiles(pg, ul.username)
		if err != nil {
			log.Printf("failed to load userfiles: %v", err)
		}
		cache.Delete(ul.username)
		cache.Add(ul.username, userfiles)
	}()

	return newEventResult(utils.SUCCESS, "login success", map[string]string{
		string(utils.ACCESS_TOKEN):  accessToken,
		string(utils.REFRESH_TOKEN): refreshToken,
	})
}
