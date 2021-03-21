package shared_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRetryFnOnConcurrencyConflict(t *testing.T) {
	Convey("Setup", t, func() {
		retryFunc := shared.RetryOnConcurrencyConflict

		Convey("Assuming the original function returns a concurrency conflict error once", func() {
			var callCounter uint8
			howOftenToFail := uint8(1)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked concurrency error"), shared.ErrConcurrencyConflict)
				}

				return nil
			}

			Convey("When RetryOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should succeed after retrying", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Assuming the original function returns a concurrency conflict error 3 times", func() {
			var callCounter uint8
			howOftenToFail := uint8(3)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked concurrency error"), shared.ErrConcurrencyConflict)
				}

				return nil
			}

			Convey("When RetryOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should fail", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)
				})
			})
		})

		Convey("Assuming the original function returns a different error", func() {
			var callCounter uint8
			howOftenToFail := uint8(1)
			originalFunc := func() error {
				callCounter++

				if callCounter <= howOftenToFail {
					return errors.Mark(errors.New("mocked technical error"), shared.ErrTechnical)
				}

				return nil
			}

			Convey("When RetryOnConcurrencyConflict is invoked with 3 maxRetries", func() {
				retries := uint8(3)

				Convey("Then it should succeed after retrying", func() {
					err := retryFunc(originalFunc, retries)
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
