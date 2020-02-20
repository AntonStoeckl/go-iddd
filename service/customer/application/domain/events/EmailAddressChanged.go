package events

import (
	"go-iddd/service/customer/application/domain/values"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressChangedAggregateName       = "Customer"
	EmailAddressChangedMetaTimestampFormat = time.RFC3339Nano
)

type EmailAddressChanged struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             Meta
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

	eventType := reflect.TypeOf(emailAddressChanged).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	fullEventName := emailAddressChangedAggregateName + eventName

	emailAddressChanged.meta = Meta{
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(EmailAddressChangedMetaTimestampFormat),
		streamVersion: streamVersion,
	}

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
		CustomerID       string `json:"customerID"`
		EmailAddress     string `json:"emailAddress"`
		ConfirmationHash string `json:"confirmationHash"`
		Meta             Meta   `json:"meta"`
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
		meta:             UnmarshalMetaFromJSON(data, streamVersion),
	}

	return emailAddressChanged
}
