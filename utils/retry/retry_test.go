package retry

import (
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/deuna/payment/domain/domainerrors"
	"github.com/stretchr/testify/assert"
)

func Test_Retry(t *testing.T) {
	t.Run("Test WithMaxAttempts", func(t *testing.T) {
		fn := func() error {
			return nil
		}
		r := NewRetry(fn)
		r = r.WithMaxAttempts(5)
		assert.Equal(t, 5, r.maxAttempts)
	})

	t.Run("Test WithSleep", func(t *testing.T) {
		fn := func() error {
			return nil
		}
		r := NewRetry(fn)
		r = r.WithSleep(time.Second)
		assert.Equal(t, time.Second, r.sleep)
	})

	t.Run("Test WithRetryableErrorType", func(t *testing.T) {
		fn := func() error {
			return nil
		}
		r := NewRetry(fn)
		r = r.WithRetryableErrorType(errors.New("test"), fmt.Errorf("test%s", "23"))
		assert.Equal(t, 2, len(r.errorsToRetry))
	})
}

func Test_Retry_Run_FailureAfterMaxAtt(t *testing.T) {
	// Arrange
	request := "test request"

	maxAtt := 3
	sleepDuration := time.Millisecond * 100

	expectedErr := domainerrors.NewInternalServerError(errors.New("test error"))

	var fnErr error

	retry := &Retry{
		errorsToRetry: map[string]struct{}{},
		maxAttempts:   maxAtt,
		sleep:         sleepDuration,
		fn: func() error {
			fnErr = func(req string) error {
				return expectedErr
			}(request)

			if fnErr != nil {
				return fnErr
			}

			return nil
		},
	}

	retry.WithRetryableErrorType(domainerrors.NewInternalServerError(nil))
	// Act
	resp, err := retry.Run()

	// Assert
	assert.Error(t, err)
	assert.Len(t, resp, maxAtt)
	assert.Equal(t, expectedErr, resp[0])
	assert.Equal(t, expectedErr, resp[1])
	assert.Equal(t, expectedErr, resp[2])
}

func Test_Retry_Run_Success(t *testing.T) {
	// Arrange
	request := "test request"
	response := "test response"

	maxAtt := 3
	sleepDuration := time.Millisecond * 100

	var fnErr error
	var fnResp string

	retry := &Retry{
		maxAttempts: maxAtt,
		sleep:       sleepDuration,
		fn: func() error {
			fnResp, fnErr = func(req string) (string, error) {
				return response, nil
			}(request)

			if fnErr != nil {
				return fnErr
			}

			return nil
		},
	}
	// Act
	resp, err := retry.Run()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Attempts())
	assert.Nil(t, resp[1])
	assert.Equal(t, response, fnResp)
}
