package serialization

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalAndUnmarshalCustomerEvents(t *testing.T) {
	customerID := values.GenerateCustomerID()
	emailAddress := values.RebuildEmailAddress("john@doe.com")
	newEmailAddress := values.RebuildEmailAddress("john.frank@doe.com")
	confirmationHash := values.GenerateConfirmationHash(emailAddress.String())
	personName := values.RebuildPersonName("John", "Doe")
	newPersonName := values.RebuildPersonName("John Frank", "Doe")
	failureReason := "wrong confirmation hash supplied"

	var myEvents es.DomainEvents
	streamVersion := uint(1)

	myEvents = append(
		myEvents,
		events.BuildCustomerRegistered(customerID, emailAddress, confirmationHash, personName, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		events.BuildCustomerEmailAddressConfirmed(customerID, emailAddress, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		events.BuildCustomerEmailAddressConfirmationFailed(customerID, emailAddress, confirmationHash, failureReason, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		events.BuildCustomerEmailAddressChanged(customerID, newEmailAddress, confirmationHash, emailAddress, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		events.BuildCustomerNameChanged(customerID, newPersonName, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		events.BuildCustomerDeleted(customerID, emailAddress, streamVersion),
	)

	for idx, event := range myEvents {
		originalEvent := event
		streamVersion := uint(idx + 1)
		eventName := originalEvent.Meta().EventName()

		Convey(fmt.Sprintf("When %s is marshaled and unmarshaled", eventName), t, func() {
			json, err := MarshalCustomerEvent(originalEvent)
			So(err, ShouldBeNil)

			unmarshaledEvent, err := UnmarshalCustomerEvent(originalEvent.Meta().EventName(), json, streamVersion)
			So(err, ShouldBeNil)

			Convey(fmt.Sprintf("Then the unmarshaled %s should resemble the original %s", eventName, eventName), func() {
				So(unmarshaledEvent, ShouldResemble, originalEvent)
			})
		})
	}
}

func TestMarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is marshaled", t, func() {
		_, err := MarshalCustomerEvent(SomeEvent{})

		Convey("Then it should fail", func() {
			So(errors.Is(err, lib.ErrMarshalingFailed), ShouldBeTrue)
		})
	})
}

func TestUnmarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		_, err := UnmarshalCustomerEvent("unknown", []byte{}, 1)

		Convey("Then it should fail", func() {
			So(errors.Is(err, lib.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}

/***** a mock event to test marshaling unknown event *****/

type SomeEvent struct{}

func (event SomeEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta("SomeEvent", "never", 1)
}

func (event SomeEvent) EventName() string                { return "SomeEvent" }
func (event SomeEvent) OccurredAt() string               { return "never" }
func (event SomeEvent) StreamVersion() uint              { return 1 }
func (event SomeEvent) IndicatesAnError() (bool, string) { return false, "" }
