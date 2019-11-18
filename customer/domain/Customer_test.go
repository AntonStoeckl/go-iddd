package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerStreamVersion(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.RebuildCustomerID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		newEmailAddress, err := values.NewEmailAddress("john+changed@doe.com")
		So(err, ShouldBeNil)
		newConfirmableEmailAddress := newEmailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(2)

		customer, err := domain.ReconstituteCustomerFrom(
			shared.DomainEvents{
				events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion),
				events.EmailAddressWasChanged(id, newConfirmableEmailAddress, currentStreamVersion),
			},
		)
		So(err, ShouldBeNil)

		Convey("When it's streamVersion is retrieved", func() {
			streamVersion := customer.StreamVersion()

			Convey("It should expose the expected version", func() {
				So(streamVersion, ShouldResemble, currentStreamVersion)
			})
		})
	})
}

func TestReconstituteCustomerFromWithInvalidEventStream(t *testing.T) {
	Convey("When a Customer is reconstituted from an empty EventStream", t, func() {
		var eventStream shared.DomainEvents

		_, err := domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})

	Convey("When a Customer is reconstituted from an EventStream without a Registered event", t, func() {
		id, err := values.RebuildCustomerID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)

		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		eventStream := shared.DomainEvents{
			events.EmailAddressWasConfirmed(id, emailAddress, uint(2)),
		}

		_, err = domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})
}
