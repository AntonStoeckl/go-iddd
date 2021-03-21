package shared

import (
	"github.com/cockroachdb/errors"
)

func RetryOnConcurrencyConflict(
	retriedFn func() error,
	maxRetries uint8,
) error {

	var err error

	for retries := uint8(0); retries < maxRetries; retries++ {
		if err = retriedFn(); err == nil {
			return nil // no need to retry, call to retriedFn was successful
		}

		if !errors.Is(err, ErrConcurrencyConflict) {
			return err // don't retry - it's not a concurrency conflice error
		}
	}

	return errors.Wrap(err, ErrMaxRetriesExceeded.Error())
}
