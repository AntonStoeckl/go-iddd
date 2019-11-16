package events

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

type Meta struct {
	identifier    string
	eventName     string
	occurredAt    string
	streamVersion uint
}

/*** Implement json.Marshaler ***/

func (meta *Meta) MarshalJSON() ([]byte, error) {
	data := &struct {
		Identifier    string `json:"identifier"`
		EventName     string `json:"eventName"`
		OccurredAt    string `json:"occurredAt"`
		StreamVersion uint   `json:"streamVersion"`
	}{
		Identifier:    meta.identifier,
		EventName:     meta.eventName,
		OccurredAt:    meta.occurredAt,
		StreamVersion: meta.streamVersion,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (meta *Meta) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		Identifier    string `json:"identifier"`
		EventName     string `json:"eventName"`
		OccurredAt    string `json:"occurredAt"`
		StreamVersion uint   `json:"streamVersion"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "meta.UnmarshalJSON")
	}

	meta.identifier = unmarshaledData.Identifier
	meta.eventName = unmarshaledData.EventName
	meta.occurredAt = unmarshaledData.OccurredAt
	meta.streamVersion = unmarshaledData.StreamVersion

	return nil
}
