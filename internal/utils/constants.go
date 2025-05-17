package utils

import (
	"errors"
	"time"
)

// tokens
type TokenType string

const (
	ACCESS_TOKEN  TokenType = "accessToken"
	REFRESH_TOKEN TokenType = "refreshToken"

	ACCESSTOKEN_MAX_AGE  = time.Minute * 20
	REFRESHTOKEN_MAX_AGE = time.Hour * 24 * 7

	USERNAME_KEY = "username"

	FILES_MAX_MEMORY = 32 << 20

	UPLOADDIR_TEMP  = "kmem_upload"
	UPLOADDIR_FINAL = "/home/kang/Downloads" // TODO: this need to be changed
)

// events
type EventStatus int

const (
	SUCCESS EventStatus = iota
	FAIL    EventStatus = iota + 1
)

// errors
var (
	TOKEN_NOT_FOUND = errors.New("access token not found")
	INVALID_TOKEN   = errors.New("invalid access token")
)

// content types
type ContentType string

const (
	IMAGE        ContentType = "image"
	VIDEO        ContentType = "video"
	AUDIO        ContentType = "audio"
	ZIP          ContentType = "zip"
	TEXT         ContentType = "text"
	OCTET_STREAM ContentType = "octet-stream"
)
