package events

import (
	"go-iddd/service/customer/application/domain/values"

	jsoniter "github.com/json-iterator/go"
)

type CustomerRegistered struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
	meta             EventMeta
}

func CustomerWasRegistered(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	personName values.PersonName,
	streamVersion uint,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
	}

	event.meta = BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerRegistered) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerRegistered) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerRegistered) ConfirmationHash() values.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerRegistered) PersonName() values.PersonName {
	return event.personName
}

func (event CustomerRegistered) EventName() string {
	return event.meta.eventName
}

func (event CustomerRegistered) OccurredAt() string {
	return event.meta.occurredAt
}

func (event CustomerRegistered) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerRegistered) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID       string    `json:"customerID"`
		EmailAddress     string    `json:"emailAddress"`
		ConfirmationHash string    `json:"confirmationHash"`
		PersonGivenName  string    `json:"personGivenName"`
		PersonFamilyName string    `json:"personFamilyName"`
		Meta             EventMeta `json:"meta"`
	}{
		CustomerID:       event.customerID.ID(),
		EmailAddress:     event.emailAddress.EmailAddress(),
		ConfirmationHash: event.confirmationHash.Hash(),
		PersonGivenName:  event.personName.GivenName(),
		PersonFamilyName: event.personName.FamilyName(),
		Meta:             event.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalCustomerRegisteredFromJSON(data []byte, streamVersion uint) CustomerRegistered {
	event := CustomerRegistered{
		customerID:       values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(jsoniter.Get(data, "confirmationHash").ToString()),
		personName: values.RebuildPersonName(
			jsoniter.Get(data, "personGivenName").ToString(),
			jsoniter.Get(data, "personFamilyName").ToString(),
		),
		meta: UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return event
}
