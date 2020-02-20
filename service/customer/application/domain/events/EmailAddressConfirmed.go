package events

import (
	"go-iddd/service/customer/application/domain/values"

	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressConfirmedAggregateName = "Customer"
)

type EmailAddressConfirmed struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         EventMeta
}

func EmailAddressWasConfirmed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
) EmailAddressConfirmed {

	emailAddressConfirmed := EmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	emailAddressConfirmed.meta = BuildEventMeta(
		emailAddressConfirmed,
		emailAddressConfirmedAggregateName,
		streamVersion,
	)

	return emailAddressConfirmed
}

func (emailAddressConfirmed EmailAddressConfirmed) CustomerID() values.CustomerID {
	return emailAddressConfirmed.customerID
}

func (emailAddressConfirmed EmailAddressConfirmed) EmailAddress() values.EmailAddress {
	return emailAddressConfirmed.emailAddress
}

func (emailAddressConfirmed EmailAddressConfirmed) EventName() string {
	return emailAddressConfirmed.meta.eventName
}

func (emailAddressConfirmed EmailAddressConfirmed) OccurredAt() string {
	return emailAddressConfirmed.meta.occurredAt
}

func (emailAddressConfirmed EmailAddressConfirmed) StreamVersion() uint {
	return emailAddressConfirmed.meta.streamVersion
}

func (emailAddressConfirmed EmailAddressConfirmed) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID   string    `json:"customerID"`
		EmailAddress string    `json:"emailAddress"`
		Meta         EventMeta `json:"meta"`
	}{
		CustomerID:   emailAddressConfirmed.customerID.ID(),
		EmailAddress: emailAddressConfirmed.emailAddress.EmailAddress(),
		Meta:         emailAddressConfirmed.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalEmailAddressConfirmedFromJSON(data []byte, streamVersion uint) EmailAddressConfirmed {
	emailAddressConfirmed := EmailAddressConfirmed{
		customerID:   values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress: values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		meta:         UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return emailAddressConfirmed
}
