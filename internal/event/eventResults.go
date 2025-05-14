package event

import "fmt"

type eventStatus int

const (
	SUCCESS eventStatus = iota
	FAIL    eventStatus = iota + 1
)

type Result struct {
	status  eventStatus
	message string
}

func newEventResult(status eventStatus, msg string) Result {
	return Result{
		status:  status,
		message: msg,
	}
}

func (r Result) String() string {
	status := "Fail"

	if r.status == SUCCESS {
		status = "Success"
	}

	return fmt.Sprintf("[%s]: %s", status, r.message)
}

func (r Result) Status() eventStatus {
	return r.status
}

func (r Result) Message() string {
	return r.message
}
