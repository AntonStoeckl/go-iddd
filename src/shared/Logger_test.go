package shared

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStandardLogger(t *testing.T) {
	Convey("When a new standard logger is created", t, func() {
		logger := NewStandardLogger()

		Convey("And it be configured to log verbose", func() {
			So(logger.Verbose(), ShouldBeTrue)
		})
	})

	Convey("When a new nil logger is created", t, func() {
		logger := NewNilLogger()

		Convey("It should be configured to log verbose", func() {
			So(logger.Verbose(), ShouldBeTrue)
		})
	})
}
