package cqrs

import (
	"go-iddd/service/lib"

	"github.com/cockroachdb/errors"
)

func RetryCommand(commandHandlerFunction func() error, maxRetries uint8) error {
	var err error
	var retries uint8

	for retries = 0; retries < maxRetries; retries++ {
		// call next method in chain
		if err = commandHandlerFunction(); err == nil {
			break // no need to retry, handling was successful
		}

		if errors.Is(err, lib.ErrConcurrencyConflict) {
			continue // retry to resolve the concurrency conflict
		} else {
			break // don't retry for different errors
		}
	}

	if err != nil {
		if retries == maxRetries {
			return errors.Wrap(err, lib.ErrMaxRetriesExceeded.Error())
		}

		return err // either to many retries or a different error
	}

	return nil
}
