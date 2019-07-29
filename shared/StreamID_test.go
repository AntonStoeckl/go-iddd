package shared_test

import (
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStreamID(t *testing.T) {
	Convey("Given valid input", t, func() {
		streamIDInput := "customer-123"

		Convey("When a new StreamID is created", func() {
			streamID := shared.NewStreamID(streamIDInput)

			Convey("It should succeed", func() {
				So(streamID, ShouldNotBeNil)
			})
		})
	})

	Convey("Given empty input", t, func() {
		streamIDInput := ""

		Convey("When a new StreamID is created", func() {
			newStreamIDWithEmptyInput := func() {
				shared.NewStreamID(streamIDInput)
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
		streamID := shared.NewStreamID(streamIDInput)

		Convey("It should expose the expected value", func() {
			So(streamID.String(), ShouldEqual, streamIDInput)
		})
	})
}

func TestStreamID_Equals(t *testing.T) {
	Convey("Given a StreamID", t, func() {
		streamID := shared.NewStreamID("customer-123")

		Convey("And given an equal StreamID", func() {
			equalStreamID := shared.NewStreamID("customer-123")

			Convey("When they are compared", func() {
				streamIDsAreEqual := streamID.Equals(equalStreamID)

				Convey("They should be equal", func() {
					So(streamIDsAreEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given a different StreamID", func() {
			differentStreamID := shared.NewStreamID("customer-666")

			Convey("When they are compared", func() {
				streamIDsAreEqual := streamID.Equals(differentStreamID)

				Convey("They should be equal", func() {
					So(streamIDsAreEqual, ShouldBeFalse)
				})
			})
		})
	})
}
