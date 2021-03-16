package serialization_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalAndUnmarshalCustomerEvents(t *testing.T) {
	customerID := value.GenerateCustomerID()
	emailAddressInput := "john@doe.com"
	confirmationHash := value.GenerateConfirmationHash(emailAddressInput)
	unconfirmedEmailAddress := value.RebuildUnconfirmedEmailAddress(emailAddressInput, confirmationHash.String())
	confirmedEmailAddress := value.RebuildConfirmedEmailAddress(emailAddressInput)
	changedEmailAddressInput := "john.frank@doe.com"
	changedConfirmationHash := value.GenerateConfirmationHash(changedEmailAddressInput)
	changedEmailAddress := value.RebuildUnconfirmedEmailAddress(changedEmailAddressInput, changedConfirmationHash.String())
	personName := value.RebuildPersonName("John", "Doe")
	newPersonName := value.RebuildPersonName("John Frank", "Doe")
	failureReason := "wrong confirmation hash supplied"
	causationID := es.GenerateMessageID()

	var myEvents []es.DomainEvent
	streamVersion := uint(1)

	myEvents = append(
		myEvents,
		domain.BuildCustomerRegistered(customerID, unconfirmedEmailAddress, personName, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerEmailAddressConfirmed(customerID, confirmedEmailAddress, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildCustomerEmailAddressChanged(customerID, changedEmailAddress, causationID, streamVersion),
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
			json, err := serialization.MarshalCustomerEvent(originalEvent)
			So(err, ShouldBeNil)

			unmarshaledEvent, err := serialization.UnmarshalCustomerEvent(originalEvent.Meta().EventName(), json, streamVersion)
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
			confirmationHash,
			errors.Mark(errors.New(failureReason), shared.ErrDomainConstraintsViolation),
			causationID,
			streamVersion,
		)

		oEventName := originalEvent.Meta().EventName()

		json, err := serialization.MarshalCustomerEvent(originalEvent)
		So(err, ShouldBeNil)

		unmarshaledEvent, err := serialization.UnmarshalCustomerEvent(originalEvent.Meta().EventName(), json, streamVersion)
		So(err, ShouldBeNil)

		uEventName := unmarshaledEvent.Meta().EventName()

		Convey(fmt.Sprintf("Then the unmarshaled %s should resemble the original %s", oEventName, uEventName), func() {
			unmarshaledEvent, ok := unmarshaledEvent.(domain.CustomerEmailAddressConfirmationFailed)
			So(ok, ShouldBeTrue)
			So(unmarshaledEvent.CustomerID().Equals(originalEvent.CustomerID()), ShouldBeTrue)
			So(unmarshaledEvent.ConfirmationHash().Equals(originalEvent.ConfirmationHash()), ShouldBeTrue)
			serialization.AssertEventMetaResembles(originalEvent, unmarshaledEvent)
		})
	})
}

func TestMarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is marshaled", t, func() {
		_, err := serialization.MarshalCustomerEvent(serialization.SomeEvent{})

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrMarshalingFailed), ShouldBeTrue)
		})
	})
}

func TestUnmarshalCustomerEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		_, err := serialization.UnmarshalCustomerEvent("unknown", []byte{}, 1)

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}
