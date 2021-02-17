package serialization

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalAndUnmarshalCustomerEvents(t *testing.T) {
	customerID := value.GenerateCustomerID()
	emailAddress := value.RebuildEmailAddress("john@doe.com")
	newEmailAddress := value.RebuildEmailAddress("john.frank@doe.com")
	confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
	personName := value.RebuildPersonName("John", "Doe")
	newPersonName := value.RebuildPersonName("John Frank", "Doe")
	failureReason := "wrong confirmation hash supplied"
	causationID := es.GenerateMessageID()

	var myEvents []es.DomainEvent
	streamVersion := uint(1)

	myEvents = append(
		myEvents,
		domain.BuildCustomerRegistered(customerID, emailAddress, confirmationHash, personName, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerEmailAddressConfirmed(customerID, emailAddress, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerEmailAddressChanged(customerID, newEmailAddress, confirmationHash, emailAddress, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerNameChanged(customerID, newPersonName, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerDeleted(customerID, causationID, streamVersion),
	)

	for idx, event := range myEvents {
		originalEvent := event
		streamVersion = uint(idx + 1)
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

	// Special treatment for Failure events because the FailureReason()
	//  is a pointer to an error which does not resemble properly (ShouldResemble uses reflect.DeepEqual)

	Convey("When CustomerEmailAddressConfirmationFailed is marshaled and unmarshaled", t, func() {
		originalEvent := domain.BuildCustomerEmailAddressConfirmationFailed(
			customerID,
			emailAddress,
			confirmationHash,
			errors.Mark(errors.New(failureReason), shared.ErrDomainConstraintsViolation),
			causationID,
			streamVersion,
		)

		oEventName := originalEvent.Meta().EventName()

		json, err := MarshalCustomerEvent(originalEvent)
		So(err, ShouldBeNil)

		unmarshaledEvent, err := UnmarshalCustomerEvent(originalEvent.Meta().EventName(), json, streamVersion)
		So(err, ShouldBeNil)

		uEventName := unmarshaledEvent.Meta().EventName()

		Convey(fmt.Sprintf("Then the unmarshaled %s should resemble the original %s", oEventName, uEventName), func() {
			unmarshaledEvent, ok := unmarshaledEvent.(domain.CustomerEmailAddressConfirmationFailed)
			So(ok, ShouldBeTrue)
			So(unmarshaledEvent.CustomerID().Equals(originalEvent.CustomerID()), ShouldBeTrue)
			So(unmarshaledEvent.EmailAddress().Equals(originalEvent.EmailAddress()), ShouldBeTrue)
			So(unmarshaledEvent.ConfirmationHash().Equals(originalEvent.ConfirmationHash()), ShouldBeTrue)
			assertEventMetaResembles(originalEvent, unmarshaledEvent)
		})
	})
}

func assertEventMetaResembles(originalEvent, unmarshaledEvent es.DomainEvent) {
	So(unmarshaledEvent.Meta().EventName(), ShouldEqual, originalEvent.Meta().EventName())
	So(unmarshaledEvent.Meta().OccurredAt(), ShouldEqual, originalEvent.Meta().OccurredAt())
	So(unmarshaledEvent.Meta().CausationID(), ShouldEqual, originalEvent.Meta().CausationID())
	So(unmarshaledEvent.Meta().StreamVersion(), ShouldEqual, originalEvent.Meta().StreamVersion())
	So(unmarshaledEvent.IsFailureEvent(), ShouldEqual, originalEvent.IsFailureEvent())
	So(unmarshaledEvent.FailureReason(), ShouldBeError)
	So(unmarshaledEvent.FailureReason().Error(), ShouldEqual, originalEvent.FailureReason().Error())
	So(errors.Is(originalEvent.FailureReason(), shared.ErrDomainConstraintsViolation), ShouldBeTrue)
	So(errors.Is(unmarshaledEvent.FailureReason(), shared.ErrDomainConstraintsViolation), ShouldBeTrue)
}

func TestMarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is marshaled", t, func() {
		_, err := MarshalCustomerEvent(SomeEvent{})

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrMarshalingFailed), ShouldBeTrue)
		})
	})
}

func TestUnmarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		_, err := UnmarshalCustomerEvent("unknown", []byte{}, 1)

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}

/***** a mock event to test marshaling unknown event *****/

type SomeEvent struct{}

func (event SomeEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta("SomeEvent", "never", "someID", "someID", 1)
}

func (event SomeEvent) IsFailureEvent() bool {
	return false
}

func (event SomeEvent) FailureReason() error {
	return nil
}
