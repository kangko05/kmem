package utils

import "time"

// tokens
const (
	ACCESS_TOKEN_DUR  = 20 * time.Minute
	REFRESH_TOKEN_DUR = (60 * 24 * 7) * time.Minute

	// key used in cookie
	ACCESS_TOKEN_KEY  = "accessToken"
	REFRESH_TOKEN_KEY = "refreshToken"
)
