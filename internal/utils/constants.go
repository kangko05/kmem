package utils

import "time"

// tokens
const (
	ACCESS_TOKEN  = "accessToken"
	REFRESH_TOKEN = "refreshToken"

	ACCESSTOKEN_MAX_AGE  = time.Minute * 20
	REFRESHTOKEN_MAX_AGE = time.Hour * 24 * 7
)

// events
type EventStatus int

const (
	SUCCESS EventStatus = iota
	FAIL    EventStatus = iota + 1
)
