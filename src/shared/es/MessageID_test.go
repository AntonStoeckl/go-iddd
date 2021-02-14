package es_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageID_Generate(t *testing.T) {
	Convey("When a MessageID is generated", t, func() {
		messageID := es.GenerateMessageID()

		Convey("It should not be empty", func() {
			So(messageID, ShouldNotBeEmpty)
		})
	})
}

func TestMessageID_Build(t *testing.T) {
	Convey("When a MessageID is built from another MessageID", t, func() {
		otherMessageID := es.GenerateMessageID()
		messageID := es.BuildMessageID(otherMessageID)

		Convey("It should not be empty", func() {
			So(messageID, ShouldNotBeEmpty)

			Convey("And it should equal the input messageID", func() {
				So(messageID.Equals(otherMessageID), ShouldBeTrue)
			})
		})
	})
}

func TestMessageID_Rebuild(t *testing.T) {
	Convey("When a MessageID is rebuilt from string", t, func() {
		messageIDString := uuid.New().String()
		messageID := es.RebuildMessageID(messageIDString)

		Convey("It should not be empty", func() {
			So(messageID, ShouldNotBeEmpty)

			Convey("And it should expose the expected value", func() {
				So(messageID.String(), ShouldEqual, messageIDString)
			})
		})
	})
}
