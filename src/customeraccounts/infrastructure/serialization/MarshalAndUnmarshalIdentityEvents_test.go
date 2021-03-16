package serialization_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalAndUnmarshalIdentityEvents(t *testing.T) {
	customerID := value.GenerateIdentityID()
	emailAddressInput := "john@doe.com"
	confirmationHash := value.GenerateConfirmationHash(emailAddressInput)
	unconfirmedEmailAddress := value.RebuildUnconfirmedEmailAddress(emailAddressInput, confirmationHash.String())
	hashedPW := value.RebuildHashedPassword("superSecurePW123")
	causationID := es.GenerateMessageID()

	var myEvents []es.DomainEvent
	streamVersion := uint(1)

	myEvents = append(
		myEvents,
		domain.BuildIdentityRegistered(customerID, unconfirmedEmailAddress, hashedPW, causationID, streamVersion),
	)

	streamVersion++

	myEvents = append(
		myEvents,
		domain.BuildIdentityDeleted(customerID, causationID, streamVersion),
	)

	for idx, event := range myEvents {
		originalEvent := event
		streamVersion = uint(idx + 1)
		eventName := originalEvent.Meta().EventName()

		Convey(fmt.Sprintf("When %s is marshaled and unmarshaled", eventName), t, func() {
			json, err := serialization.MarshalIdentityEvent(originalEvent)
			So(err, ShouldBeNil)

			unmarshaledEvent, err := serialization.UnmarshalIdentityEvent(originalEvent.Meta().EventName(), json, streamVersion)
			So(err, ShouldBeNil)

			Convey(fmt.Sprintf("Then the unmarshaled %s should resemble the original %s", eventName, eventName), func() {
				So(unmarshaledEvent, ShouldResemble, originalEvent)
			})
		})
	}
}

func TestMarshalIdentityEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is marshaled", t, func() {
		_, err := serialization.MarshalIdentityEvent(serialization.SomeEvent{})

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrMarshalingFailed), ShouldBeTrue)
		})
	})
}

func TestUnmarshalIdentityEvent_WithUnknownEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		_, err := serialization.UnmarshalIdentityEvent("unknown", []byte{}, 1)

		Convey("Then it should fail", func() {
			So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}
