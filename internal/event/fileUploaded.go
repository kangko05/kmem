package event

import (
	"context"
	"kmem/internal/utils"
	"time"
)

/*
	handle(context.Context) Result

	setResultChannel(chan Result)
	getResultChannel() chan Result

	setTimeout(time.Duration)
*/

type fileUploaded struct {
	resultCh chan Result
	timeout  time.Duration
}

func FileUploaded(options ...eventOption) *fileUploaded {
	fileUploaded := &fileUploaded{
		resultCh: nil,
		timeout:  time.Duration(defaultTimeout),
	}

	return fileUploaded
}

func (f *fileUploaded) setResultChannel(resultChan chan Result) {
	f.resultCh = resultChan
}

func (f *fileUploaded) getResultChannel() chan Result {
	return f.resultCh
}

func (f *fileUploaded) setTimeout(dur time.Duration) {
	f.timeout = dur
}

func (f *fileUploaded) handle(ctx context.Context) Result {
	// 1. gather meta data
	// 2. meta data into db
	// 3. store files into local path (truenas later)

	return newEventResult(utils.FAIL, "", nil)
}
