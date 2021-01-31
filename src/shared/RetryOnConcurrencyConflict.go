package shared

import (
	"github.com/cockroachdb/errors"
)

func RetryOnConcurrencyConflict(originalFunc func() error, maxRetries uint8) error {
	var err error
	var retries uint8

	for retries = 0; retries < maxRetries; retries++ {
		// call next method in chain
		if err = originalFunc(); err == nil {
			return nil // no need to retry, call to originalFunc was successful
		}

		if !errors.Is(err, ErrConcurrencyConflict) {
			return err // don't retry for different errors
		}
	}

	return errors.Wrap(err, ErrMaxRetriesExceeded.Error())
}
