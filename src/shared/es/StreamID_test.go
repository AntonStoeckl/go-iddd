package es_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStreamID_BuildStreamID(t *testing.T) {
	Convey("Given valid input", t, func() {
		streamIDInput := "customer-123"

		Convey("When a new StreamID is created", func() {
			streamID := es.BuildStreamID(streamIDInput)

			Convey("It should succeed", func() {
				So(streamID, ShouldNotBeNil)
			})
		})
	})

	Convey("Given empty input", t, func() {
		streamIDInput := ""

		Convey("When a new StreamID is created", func() {
			newStreamIDWithEmptyInput := func() {
				es.BuildStreamID(streamIDInput)
			}

			Convey("It should fail with a panic", func() {
				So(newStreamIDWithEmptyInput, ShouldPanic)
			})
		})
	})
}

func TestStreamID_String(t *testing.T) {
	Convey("Given a StreamID", t, func() {
		streamIDInput := "customer-123"
		streamID := es.BuildStreamID(streamIDInput)

		Convey("It should expose the expected value", func() {
			So(streamID.String(), ShouldEqual, streamIDInput)
		})
	})
}
