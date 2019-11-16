package events_test

import (
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkUnmarshalDomainEvent(b *testing.B) {
	id := values.GenerateID()
	emailAddress, _ := values.NewEmailAddress("foo@bar.com")
	confirmableEmailAddress := emailAddress.ToConfirmable()
	personName, _ := values.NewPersonName("John", "Doe")
	streamVersion := uint(1)

	event1 := events.ItWasRegistered(id, confirmableEmailAddress, personName, streamVersion)
	data1, _ := event1.MarshalJSON()

	event2 := events.EmailAddressWasConfirmed(id, emailAddress, streamVersion)
	data2, _ := event2.MarshalJSON()

	event3 := events.EmailAddressWasChanged(id, confirmableEmailAddress, streamVersion)
	data3, _ := event3.MarshalJSON()

	for n := 0; n < b.N; n++ {
		_, _ = events.UnmarshalDomainEvent("CustomerRegistered", data1)
		_, _ = events.UnmarshalDomainEvent("CustomerEmailAddressConfirmed", data2)
		_, _ = events.UnmarshalDomainEvent("CustomerEmailAddressChanged", data3)
	}
}

func TestUnmarshalRegistered(t *testing.T) {
	Convey("When CustomerRegistered is unmarshaled", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)
		streamVersion := uint(1)

		Convey("And when the input is valid json", func() {
			event := events.ItWasRegistered(id, confirmableEmailAddress, personName, streamVersion)
			data, err := event.MarshalJSON()
			So(err, ShouldBeNil)

			unmarshaled, err := events.UnmarshalDomainEvent("CustomerRegistered", data)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(event, ShouldResemble, unmarshaled)
			})
		})

		Convey("And when the input is invalid json", func() {
			unmarshaled, err := events.UnmarshalDomainEvent("CustomerRegistered", []byte("666"))

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(unmarshaled, ShouldBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}

func TestUnmarshalCustomerEmailAddressConfirmed(t *testing.T) {
	Convey("When CustomerEmailAddressConfirmed is unmarshaled", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		streamVersion := uint(1)

		Convey("And when the input is valid json", func() {
			event := events.EmailAddressWasConfirmed(id, emailAddress, streamVersion)
			data, err := event.MarshalJSON()
			So(err, ShouldBeNil)

			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressConfirmed", data)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(event, ShouldResemble, unmarshaled)
			})
		})

		Convey("And when the input is invalid json", func() {
			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressConfirmed", []byte("666"))

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(unmarshaled, ShouldBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}

func TestUnmarshalEmailAddressConfirmationFailed(t *testing.T) {
	Convey("When EmailAddressConfirmationFailed is unmarshaled", t, func() {
		id := values.GenerateID()
		streamVersion := uint(1)
		invalidHash := values.GenerateConfirmationHash("invalid_hash")

		Convey("And when the input is valid json", func() {
			event := events.EmailAddressConfirmationHasFailed(id, invalidHash, streamVersion)
			data, err := event.MarshalJSON()
			So(err, ShouldBeNil)

			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressConfirmationFailed", data)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(event, ShouldResemble, unmarshaled)
			})
		})

		Convey("And when the input is invalid json", func() {
			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressConfirmationFailed", []byte("666"))

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(unmarshaled, ShouldBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}

func TestUnmarshalCustomerEmailAddressChanged(t *testing.T) {
	Convey("When CustomerEmailAddressChanged is unmarshaled", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		streamVersion := uint(1)

		Convey("And when the input is valid json", func() {
			event := events.EmailAddressWasChanged(id, confirmableEmailAddress, streamVersion)
			data, err := event.MarshalJSON()
			So(err, ShouldBeNil)

			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressChanged", data)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(event, ShouldResemble, unmarshaled)
			})
		})

		Convey("And when the input is invalid json", func() {
			unmarshaled, err := events.UnmarshalDomainEvent("CustomerEmailAddressChanged", []byte("666"))

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(unmarshaled, ShouldBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}

func TestUnmarshalWithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		unmarshaled, err := events.UnmarshalDomainEvent("UnknownEvent", []byte("666"))

		Convey("It should fail", func() {
			So(err, ShouldNotBeNil)
			So(unmarshaled, ShouldBeNil)
			So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}
