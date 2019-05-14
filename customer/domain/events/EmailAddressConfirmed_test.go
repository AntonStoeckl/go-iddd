package events_test

import (
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestEmailAddressWasConfirmed(t *testing.T) {
	Convey("Given valid parameters as input", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		Convey("When a new EmailAddressConfirmed event is created", func() {
			emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)

			Convey("It should succeed", func() {
				So(emailAddressConfirmed, ShouldNotBeNil)
				So(emailAddressConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))
			})
		})
	})
}

func TestEmailAddressConfirmedExposesExpectedValues(t *testing.T) {
	Convey("Given a EmailAddressConfirmed event", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		beforeItOccurred := time.Now()
		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		afterItOccurred := time.Now()
		So(emailAddressConfirmed, ShouldNotBeNil)
		So(emailAddressConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		Convey("It should expose the expected values", func() {
			So(emailAddressConfirmed.ID(), ShouldResemble, id)
			So(emailAddressConfirmed.EmailAddress(), ShouldResemble, emailAddress)
			So(emailAddressConfirmed.Identifier(), ShouldEqual, id.String())
			So(emailAddressConfirmed.EventName(), ShouldEqual, "CustomerEmailAddressConfirmed")
			itOccurred, err := time.Parse(shared.DomainEventMetaTimestampFormat, emailAddressConfirmed.OccurredAt())
			So(err, ShouldBeNil)
			So(beforeItOccurred, ShouldHappenBefore, itOccurred)
			So(afterItOccurred, ShouldHappenAfter, itOccurred)
		})
	})
}

func TestEmailAddressConfirmedMarshalJSON(t *testing.T) {
	Convey("Given a EmailAddressConfirmed event", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		So(emailAddressConfirmed, ShouldNotBeNil)
		So(emailAddressConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := emailAddressConfirmed.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})
		})
	})
}

func TestEmailAddressConfirmedUnmarshalJSON(t *testing.T) {
	Convey("Given a EmailAddressConfirmed event marshaled to json", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress)
		So(emailAddressConfirmed, ShouldNotBeNil)
		So(emailAddressConfirmed, ShouldHaveSameTypeAs, (*events.EmailAddressConfirmed)(nil))

		data, err := emailAddressConfirmed.MarshalJSON()

		Convey("And when it is unmarshaled", func() {
			unmarshaled := &events.EmailAddressConfirmed{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original EmailAddressConfirmed event", func() {
				So(err, ShouldBeNil)
				So(emailAddressConfirmed, ShouldResemble, unmarshaled)
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