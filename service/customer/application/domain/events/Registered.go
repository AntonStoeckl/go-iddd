package events

import (
	"go-iddd/service/customer/application/domain/values"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	registeredAggregateName       = "Customer"
	RegisteredMetaTimestampFormat = time.RFC3339Nano
)

type Registered struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
	meta             Meta
}

func ItWasRegistered(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	personName values.PersonName,
	streamVersion uint,
) Registered {

	registered := Registered{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
	}

	eventType := reflect.TypeOf(registered).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	fullEventName := registeredAggregateName + eventName

	registered.meta = Meta{
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(RegisteredMetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return registered
}

func (registered Registered) CustomerID() values.CustomerID {
	return registered.customerID
}

func (registered Registered) EmailAddress() values.EmailAddress {
	return registered.emailAddress
}

func (registered Registered) ConfirmationHash() values.ConfirmationHash {
	return registered.confirmationHash
}

func (registered Registered) PersonName() values.PersonName {
	return registered.personName
}

func (registered Registered) EventName() string {
	return registered.meta.eventName
}

func (registered Registered) OccurredAt() string {
	return registered.meta.occurredAt
}

func (registered Registered) StreamVersion() uint {
	return registered.meta.streamVersion
}

func (registered Registered) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID       string `json:"customerID"`
		EmailAddress     string `json:"emailAddress"`
		ConfirmationHash string `json:"confirmationHash"`
		PersonGivenName  string `json:"personGivenName"`
		PersonFamilyName string `json:"personFamilyName"`
		Meta             Meta   `json:"meta"`
	}{
		CustomerID:       registered.customerID.ID(),
		EmailAddress:     registered.emailAddress.EmailAddress(),
		ConfirmationHash: registered.confirmationHash.Hash(),
		PersonGivenName:  registered.personName.GivenName(),
		PersonFamilyName: registered.personName.FamilyName(),
		Meta:             registered.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalRegisteredFromJSON(data []byte, streamVersion uint) Registered {
	registered := Registered{
		customerID:       values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(jsoniter.Get(data, "confirmationHash").ToString()),
		personName: values.RebuildPersonName(
			jsoniter.Get(data, "personGivenName").ToString(),
			jsoniter.Get(data, "personFamilyName").ToString(),
		),
		meta: UnmarshalMetaFromJSON(data, streamVersion),
	}

	return registered
}
