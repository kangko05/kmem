package event

import (
	"context"
	"fmt"
	"time"
)

// interface
type eventType interface {
	handle(context.Context) Result

	setResultChannel(chan Result)
	getResultChannel() chan Result

	setTimeout(time.Duration)
}

// options

type eventOption func(eventType)

func WithResultChan(rchan chan Result) func(eventType) {
	return func(et eventType) {
		et.setResultChannel(rchan)
	}
}

func WithTimeout(timeout time.Duration) func(eventType) {
	return func(et eventType) {
		if timeout > 0 {
			et.setTimeout(timeout)
		}
	}
}

// test event
type testEvent struct {
	resultChan chan Result
	timeout    time.Duration
	idx        int
}

func TestEvent(idx int, options ...eventOption) *testEvent {
	te := &testEvent{resultChan: nil, idx: idx, timeout: 10 * time.Second}

	for _, opt := range options {
		opt(te)
	}

	return te
}

func (te *testEvent) getResultChannel() chan Result {
	return te.resultChan
}

func (te *testEvent) handle(ctx context.Context) Result {
	select {
	case <-ctx.Done():
		return newEventResult(FAIL, fmt.Sprintf("test event canceld: %v", ctx.Err()))
	case <-time.After(te.timeout):
		return newEventResult(FAIL, "event has timed out")
	default:
		// do something  here
		return newEventResult(SUCCESS, fmt.Sprintf("test event %d", te.idx))
	}

}

func (te *testEvent) setResultChannel(rchan chan Result) {
	te.resultChan = rchan
}

func (te *testEvent) setTimeout(t time.Duration) {
	te.timeout = t
}
