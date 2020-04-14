package cmd

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStandardLogger(t *testing.T) {
	Convey("When a new standard logger is created", t, func() {
		logger := NewStandardLogger()

		Convey("Then it should be configured to use a test formatter", func() {
			var foo *logrus.TextFormatter
			So(logger.Formatter, ShouldHaveSameTypeAs, foo)

			Convey("And it be configured to log verbose", func() {
				So(logger.Verbose(), ShouldBeTrue)
			})
		})
	})

	Convey("When a new nil logger is created", t, func() {
		logger := NewNilLogger()

		Convey("Then it should be configured to discard any log output", func() {
			So(logger.Out, ShouldEqual, ioutil.Discard)
		})
	})
}
