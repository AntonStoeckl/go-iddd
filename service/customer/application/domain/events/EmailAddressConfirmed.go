package events

import (
	"go-iddd/service/customer/application/domain/values"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressConfirmedAggregateName       = "Customer"
	EmailAddressConfirmedMetaTimestampFormat = time.RFC3339Nano
)

type EmailAddressConfirmed struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         Meta
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

	eventType := reflect.TypeOf(emailAddressConfirmed).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	fullEventName := emailAddressConfirmedAggregateName + eventName

	emailAddressConfirmed.meta = Meta{
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(EmailAddressConfirmedMetaTimestampFormat),
		streamVersion: streamVersion,
	}

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
		CustomerID   string `json:"customerID"`
		EmailAddress string `json:"emailAddress"`
		Meta         Meta   `json:"meta"`
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
		meta:         UnmarshalMetaFromJSON(data, streamVersion),
	}

	return emailAddressConfirmed
}
