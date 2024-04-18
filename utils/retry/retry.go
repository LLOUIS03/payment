package retry

import (
	"reflect"
	"time"

	"emperror.dev/errors"
)

type Retry struct {
	maxAttempts   int
	sleep         time.Duration
	fn            func() error
	errorsToRetry map[string]struct{}
}

type RetryResponse map[int]error

func (r RetryResponse) Attempts() int {
	return len(r)
}

func NewRetry(fn func() error) *Retry {
	return &Retry{
		fn:            fn,
		maxAttempts:   1,
		sleep:         0,
		errorsToRetry: map[string]struct{}{},
	}
}

func (r *Retry) WithMaxAttempts(attempts int) *Retry {
	r.maxAttempts = attempts
	return r
}

func (r *Retry) WithSleep(sleep time.Duration) *Retry {
	r.sleep = sleep
	return r
}

func (r *Retry) WithRetryableErrorType(types ...any) *Retry {
	for _, t := range types {
		r.errorsToRetry[reflect.TypeOf(t).String()] = struct{}{}
	}
	return r
}

// Run runs the function with the configured max attempts and sleep duration
func (r *Retry) Run() (RetryResponse, error) {
	var err error
	resp := make(RetryResponse)
	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		resp[attempt] = nil
		err = r.fn()
		if err == nil {
			return resp, nil
		}

		resp[attempt] = err

		if !r.shouldRetry(err) {
			break
		}

		time.Sleep(r.sleep)
	}

	return resp, err
}

func (r *Retry) shouldRetry(err error) bool {
	if len(r.errorsToRetry) == 0 {
		return false
	}

	t := reflect.TypeOf(errors.Cause(err)).String()
	_, needsRetry := r.errorsToRetry[t]
	return needsRetry
}
