package cqrs

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

func RetryCommandOnConcurrencyConflict(originalFunc func() error, maxRetries uint8) error {
	var err error
	var retries uint8

	for retries = 0; retries < maxRetries; retries++ {
		// call next method in chain
		if err = originalFunc(); err == nil {
			return nil // no need to retry, call to originalFunc was successful
		}

		if !errors.Is(err, lib.ErrConcurrencyConflict) {
			return err // don't retry for different errors
		}
	}

	return errors.Wrap(err, lib.ErrMaxRetriesExceeded.Error())
}
