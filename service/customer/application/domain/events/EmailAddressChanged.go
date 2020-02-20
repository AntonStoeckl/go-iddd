package events

import (
	"go-iddd/service/customer/application/domain/values"

	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressChangedAggregateName = "Customer"
)

type EmailAddressChanged struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             EventMeta
}

func EmailAddressWasChanged(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	streamVersion uint,
) EmailAddressChanged {

	emailAddressChanged := EmailAddressChanged{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	emailAddressChanged.meta = BuildEventMeta(
		emailAddressChanged,
		emailAddressChangedAggregateName,
		streamVersion,
	)

	return emailAddressChanged
}

func (emailAddressChanged EmailAddressChanged) CustomerID() values.CustomerID {
	return emailAddressChanged.customerID
}

func (emailAddressChanged EmailAddressChanged) EmailAddress() values.EmailAddress {
	return emailAddressChanged.emailAddress
}

func (emailAddressChanged EmailAddressChanged) ConfirmationHash() values.ConfirmationHash {
	return emailAddressChanged.confirmationHash
}

func (emailAddressChanged EmailAddressChanged) EventName() string {
	return emailAddressChanged.meta.eventName
}

func (emailAddressChanged EmailAddressChanged) OccurredAt() string {
	return emailAddressChanged.meta.occurredAt
}

func (emailAddressChanged EmailAddressChanged) StreamVersion() uint {
	return emailAddressChanged.meta.streamVersion
}

func (emailAddressChanged EmailAddressChanged) MarshalJSON() ([]byte, error) {
	data := struct {
		CustomerID       string    `json:"customerID"`
		EmailAddress     string    `json:"emailAddress"`
		ConfirmationHash string    `json:"confirmationHash"`
		Meta             EventMeta `json:"meta"`
	}{
		CustomerID:       emailAddressChanged.customerID.ID(),
		EmailAddress:     emailAddressChanged.emailAddress.EmailAddress(),
		ConfirmationHash: emailAddressChanged.confirmationHash.Hash(),
		Meta:             emailAddressChanged.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalEmailAddressChangedFromJSON(data []byte, streamVersion uint) EmailAddressChanged {
	emailAddressChanged := EmailAddressChanged{
		customerID:       values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(jsoniter.Get(data, "confirmationHash").ToString()),
		meta:             UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return emailAddressChanged
}
