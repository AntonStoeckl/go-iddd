package events

import (
	"go-iddd/service/customer/application/domain/values"

	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressConfirmationFailedAggregateName = "Customer"
)

type EmailAddressConfirmationFailed struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             EventMeta
}

func EmailAddressConfirmationHasFailed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	streamVersion uint,
) EmailAddressConfirmationFailed {

	emailAddressConfirmationFailed := EmailAddressConfirmationFailed{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	emailAddressConfirmationFailed.meta = BuildEventMeta(
		emailAddressConfirmationFailed,
		emailAddressConfirmationFailedAggregateName,
		streamVersion,
	)

	return emailAddressConfirmationFailed
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) CustomerID() values.CustomerID {
	return emailAddressConfirmationFailed.customerID
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) EmailAddress() values.EmailAddress {
	return emailAddressConfirmationFailed.emailAddress
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) ConfirmationHash() values.ConfirmationHash {
	return emailAddressConfirmationFailed.confirmationHash
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) EventName() string {
	return emailAddressConfirmationFailed.meta.eventName
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) OccurredAt() string {
	return emailAddressConfirmationFailed.meta.occurredAt
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) StreamVersion() uint {
	return emailAddressConfirmationFailed.meta.streamVersion
}

func (emailAddressConfirmationFailed EmailAddressConfirmationFailed) MarshalJSON() ([]byte, error) {
	data := struct {
		CustomerID       string    `json:"customerID"`
		EmailAddress     string    `json:"emailAddress"`
		ConfirmationHash string    `json:"confirmationHash"`
		Meta             EventMeta `json:"meta"`
	}{
		CustomerID:       emailAddressConfirmationFailed.customerID.ID(),
		EmailAddress:     emailAddressConfirmationFailed.emailAddress.EmailAddress(),
		ConfirmationHash: emailAddressConfirmationFailed.confirmationHash.Hash(),
		Meta:             emailAddressConfirmationFailed.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalEmailAddressConfirmationFailedFromJSON(data []byte, streamVersion uint) EmailAddressConfirmationFailed {
	emailAddressConfirmationFailed := EmailAddressConfirmationFailed{
		customerID:       values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(jsoniter.Get(data, "confirmationHash").ToString()),
		meta:             UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return emailAddressConfirmationFailed
}
