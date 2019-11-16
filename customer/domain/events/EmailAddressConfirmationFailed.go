// Code generated by generate/main.go. DO NOT EDIT.

package events

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

const (
	emailAddressConfirmationFailedAggregateName       = "Customer"
	EmailAddressConfirmationFailedMetaTimestampFormat = time.RFC3339Nano
)

type EmailAddressConfirmationFailed struct {
	id               *values.ID
	confirmationHash *values.ConfirmationHash
	meta             *Meta
}

/*** Factory Methods ***/

func EmailAddressConfirmationHasFailed(
	id *values.ID,
	confirmationHash *values.ConfirmationHash,
	streamVersion uint,
) *EmailAddressConfirmationFailed {

	emailAddressConfirmationFailed := &EmailAddressConfirmationFailed{
		id:               id,
		confirmationHash: confirmationHash,
	}

	eventType := reflect.TypeOf(emailAddressConfirmationFailed).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := emailAddressConfirmationFailedAggregateName + eventName

	emailAddressConfirmationFailed.meta = &Meta{
		identifier:    id.String(),
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(EmailAddressConfirmationFailedMetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return emailAddressConfirmationFailed
}

/*** Getter Methods ***/

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) ID() *values.ID {
	return emailAddressConfirmationFailed.id
}

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) ConfirmationHash() *values.ConfirmationHash {
	return emailAddressConfirmationFailed.confirmationHash
}

/*** Implement shared.DomainEvent ***/

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) Identifier() string {
	return emailAddressConfirmationFailed.meta.identifier
}

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) EventName() string {
	return emailAddressConfirmationFailed.meta.eventName
}

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) OccurredAt() string {
	return emailAddressConfirmationFailed.meta.occurredAt
}

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) StreamVersion() uint {
	return emailAddressConfirmationFailed.meta.streamVersion
}

/*** Implement json.Marshaler ***/

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID               *values.ID               `json:"id"`
		ConfirmationHash *values.ConfirmationHash `json:"confirmationHash"`
		Meta             *Meta                    `json:"meta"`
	}{
		ID:               emailAddressConfirmationFailed.id,
		ConfirmationHash: emailAddressConfirmationFailed.confirmationHash,
		Meta:             emailAddressConfirmationFailed.meta,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (emailAddressConfirmationFailed *EmailAddressConfirmationFailed) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		ID               *values.ID               `json:"id"`
		ConfirmationHash *values.ConfirmationHash `json:"confirmationHash"`
		Meta             *Meta                    `json:"meta"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "emailAddressConfirmationFailed.UnmarshalJSON")
	}

	emailAddressConfirmationFailed.id = unmarshaledData.ID
	emailAddressConfirmationFailed.confirmationHash = unmarshaledData.ConfirmationHash
	emailAddressConfirmationFailed.meta = unmarshaledData.Meta

	return nil
}
