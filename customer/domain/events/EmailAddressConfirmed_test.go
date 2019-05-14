package events_test

import (
	"encoding/json"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"
	"time"

	"golang.org/x/xerrors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailAddressWasConfirmed(t *testing.T) {
	Convey("Given valid parameters as input", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		Convey("When a new EmailAddressConfirmed event is created", func() {
			emailAddressWasConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)

			Convey("It should succeed", func() {
				So(emailAddressWasConfirmed, ShouldNotBeNil)
				So(emailAddressWasConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))
			})
		})
	})
}

func TestEmailAddressConfirmedExposesExpectedValues(t *testing.T) {
	Convey("Given an EmailAddressConfirmed event", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		beforeItOccurred := time.Now()
		emailAddressWasConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		afterItOccurred := time.Now()
		So(emailAddressWasConfirmed, ShouldNotBeNil)
		So(emailAddressWasConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		Convey("It should expose the expected values", func() {
			So(emailAddressWasConfirmed.ID(), ShouldResemble, id)
			So(emailAddressWasConfirmed.EmailAddress(), ShouldResemble, emailAddress)
			So(emailAddressWasConfirmed.Identifier(), ShouldEqual, id.String())
			So(emailAddressWasConfirmed.EventName(), ShouldEqual, "CustomerEmailAddressConfirmed")
			actualOccurredAt, err := time.Parse(time.RFC3339Nano, emailAddressWasConfirmed.OccurredAt())
			So(err, ShouldBeNil)
			So(beforeItOccurred, ShouldHappenBefore, actualOccurredAt)
			So(afterItOccurred, ShouldHappenAfter, actualOccurredAt)
		})
	})
}

func TestEmailAddressConfirmedMarshalJSON(t *testing.T) {
	Convey("Given a valid EmailAddressConfirmed event", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		emailAddressWasConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		So(emailAddressWasConfirmed, ShouldNotBeNil)
		So(emailAddressWasConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := json.Marshal(emailAddressWasConfirmed)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})

			Convey("And when it is unmarshaled from json", func() {
				unmarshaledEvent := &events.EmailAddressConfirmed{}
				err := json.Unmarshal(data, unmarshaledEvent)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(unmarshaledEvent, ShouldNotBeNil)
					So(unmarshaledEvent, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))
					So(emailAddressWasConfirmed, ShouldResemble, unmarshaledEvent)
				})
			})
		})
	})
}

func TestEmailAddressConfirmedUnmarshalJSON(t *testing.T) {
	Convey("Given an EmailAddressConfirmed event marshaled to json", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		emailAddressWasConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		So(emailAddressWasConfirmed, ShouldNotBeNil)
		So(emailAddressWasConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		data, err := emailAddressWasConfirmed.MarshalJSON()

		Convey("And when it is unmarshaled", func() {
			unmarshaled := &events.EmailAddressConfirmed{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original EmailAddressConfirmed event", func() {
				So(err, ShouldBeNil)
				So(emailAddressWasConfirmed, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to EmailAddressConfirmed event", func() {
			unmarshaled := &events.EmailAddressConfirmed{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
