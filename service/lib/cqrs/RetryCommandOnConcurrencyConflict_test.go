package cqrs_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/cqrs"
	"github.com/cockroachdb/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRetryCommandOnConcurrencyConflict(t *testing.T) {
	Convey("Setup", t, func() {
		retryFunc := cqrs.RetryCommandOnConcurrencyConflict

		Convey("Assuming the original function returns a concurrencyConflict error once", func() {
			var callCounter uint8
			howOftenToFail := uint8(1)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked concurrency error"), lib.ErrConcurrencyConflict)
				}

				return nil
			}

			Convey("When RetryCommandOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should succeed after retrying", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Assuming the original function returns a concurrencyConflict error 3 times", func() {
			var callCounter uint8
			howOftenToFail := uint8(3)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked concurrency error"), lib.ErrConcurrencyConflict)
				}

				return nil
			}

			Convey("When RetryCommandOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should fail", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
				})
			})
		})

		Convey("Assuming the original function returns a different error", func() {
			var callCounter uint8
			howOftenToFail := uint8(1)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked technical error"), lib.ErrTechnical)
				}

				return nil
			}

			Convey("When RetryCommandOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should succeed after retrying", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
