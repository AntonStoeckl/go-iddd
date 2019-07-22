package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestUnmarshalDomainEvent(t *testing.T) {
	Convey("Given valid input data for events", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)
		streamVersion := uint(1)

		Convey("When a CustomerRegistered event is unmarshaled", func() {
			Convey("And when the input is valid json", func() {
				event := events.ItWasRegistered(id, confirmableEmailAddress, personName, streamVersion)
				data, err := event.MarshalJSON()
				So(err, ShouldBeNil)

				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerRegistered", data)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(event, ShouldResemble, unmarshaled)
				})
			})

			Convey("And when the input is invalid json", func() {
				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerRegistered", []byte("666"))

				Convey("It should fail", func() {
					So(err, ShouldNotBeNil)
					So(unmarshaled, ShouldBeNil)
					So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("When a CustomerEmailAddressConfirmed event is unmarshaled", func() {
			Convey("And when the input is valid json", func() {
				event := events.EmailAddressWasConfirmed(id, emailAddress, streamVersion)
				data, err := event.MarshalJSON()
				So(err, ShouldBeNil)

				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerEmailAddressConfirmed", data)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(event, ShouldResemble, unmarshaled)
				})
			})

			Convey("And when the input is invalid json", func() {
				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerEmailAddressConfirmed", []byte("666"))

				Convey("It should fail", func() {
					So(err, ShouldNotBeNil)
					So(unmarshaled, ShouldBeNil)
					So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("When a CustomerEmailAddressChanged event is unmarshaled", func() {
			Convey("And when the input is valid json", func() {
				event := events.EmailAddressWasChanged(id, emailAddress, streamVersion)
				data, err := event.MarshalJSON()
				So(err, ShouldBeNil)

				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerEmailAddressChanged", data)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(event, ShouldResemble, unmarshaled)
				})
			})

			Convey("And when the input is invalid json", func() {
				unmarshaled, err := domain.UnmarshalDomainEvent("CustomerEmailAddressChanged", []byte("666"))

				Convey("It should fail", func() {
					So(err, ShouldNotBeNil)
					So(unmarshaled, ShouldBeNil)
					So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("When an unknown event is unmarshaled", func() {
			unmarshaled, err := domain.UnmarshalDomainEvent("UnknownEvent", []byte("666"))

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(unmarshaled, ShouldBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
