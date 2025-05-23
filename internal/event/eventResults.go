package event

import (
	"fmt"
	"kmem/internal/utils"
)

type Result struct {
	status  utils.EventStatus
	message string
	payload any
}

func newEventResult(status utils.EventStatus, msg string, payload any) Result {
	return Result{
		status:  status,
		message: msg,
		payload: payload,
	}
}

func (r Result) String() string {
	status := "Fail"

	if r.status == utils.SUCCESS {
		status = "Success"
	}

	return fmt.Sprintf("[%s]: %s", status, r.message)
}

func (r Result) Status() utils.EventStatus {
	return r.status
}

func (r Result) Message() string {
	return r.message
}

func (r Result) Payload() any {
	return r.payload
}
