// +build test

package test

import (
	"database/sql"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/eventstore"
	"time"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func SetUpDIContainer() *DIContainer {
	config, err := NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	err = db.Ping()
	So(err, ShouldBeNil)

	migrator, err := eventstore.NewMigrator(db, config.Postgres.MigrationsPath)
	So(err, ShouldBeNil)

	err = migrator.Up()
	So(err, ShouldBeNil)

	diContainer, err := NewDIContainer(db, Unmarshal)
	So(err, ShouldBeNil)

	return diContainer
}

func BeginTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	So(err, ShouldBeNil)

	return tx
}

/*** mocked CustomerID ***/

type SomeID struct {
	ID string
}

func (someID *SomeID) String() string {
	return someID.ID
}

func (someID *SomeID) Equals(other shared.IdentifiesAggregates) bool {
	return true // not needed in scope of this test
}

func (someID *SomeID) MarshalJSON() ([]byte, error) {
	bytes, err := jsoniter.Marshal(someID.ID)
	if err != nil {
		return bytes, xerrors.Errorf("SomeID.MarshalJSON: %s: %w", err, shared.ErrMarshalingFailed)
	}

	return bytes, nil
}

func (someID *SomeID) UnmarshalJSON(data []byte) error {
	var value string

	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return xerrors.Errorf("SomeID.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	someID.ID = value

	return nil
}

/*** mocked Event that works ***/

type SomeEvent struct {
	id         *SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateSomeEvent(forId *SomeID, withVersion uint) *SomeEvent {
	return &SomeEvent{
		id:         forId,
		name:       "SomeEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (someEvent *SomeEvent) Identifier() string {
	return someEvent.id.String()
}

func (someEvent *SomeEvent) EventName() string {
	return someEvent.name
}

func (someEvent *SomeEvent) OccurredAt() string {
	return someEvent.occurredAt
}

func (someEvent *SomeEvent) StreamVersion() uint {
	return someEvent.version
}

func (someEvent *SomeEvent) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID         *SomeID `json:"CustomerID"`
		Name       string  `json:"name"`
		Version    uint    `json:"version"`
		OccurredAt string  `json:"occurredAt"`
	}{
		ID:         someEvent.id,
		Name:       someEvent.name,
		Version:    someEvent.version,
		OccurredAt: someEvent.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func (someEvent *SomeEvent) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		ID         *SomeID `json:"CustomerID"`
		Name       string  `json:"name"`
		Version    uint    `json:"version"`
		OccurredAt string  `json:"occurredAt"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return xerrors.Errorf("SomeEvent.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	someEvent.id = unmarshaledData.ID
	someEvent.name = unmarshaledData.Name
	someEvent.version = unmarshaledData.Version
	someEvent.occurredAt = unmarshaledData.OccurredAt

	return nil
}

/*** mocked Event with broken marshaling ***/

type BrokenMarshalingEvent struct {
	id         *SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenMarshalingEvent(forId *SomeID, withVersion uint) *BrokenMarshalingEvent {
	return &BrokenMarshalingEvent{
		id:         forId,
		name:       "BrokenMarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) Identifier() string {
	return brokenMarshalingEvent.id.String()
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) EventName() string {
	return brokenMarshalingEvent.name
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) OccurredAt() string {
	return brokenMarshalingEvent.occurredAt
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) StreamVersion() uint {
	return brokenMarshalingEvent.version
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) MarshalJSON() ([]byte, error) {
	return nil, errors.New("mocked marshaling error")
}

func (brokenMarshalingEvent *BrokenMarshalingEvent) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		ID         *SomeID `json:"CustomerID"`
		Name       string  `json:"name"`
		Version    uint    `json:"version"`
		OccurredAt string  `json:"occurredAt"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return xerrors.Errorf("SomeEvent.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	brokenMarshalingEvent.id = unmarshaledData.ID
	brokenMarshalingEvent.name = unmarshaledData.Name
	brokenMarshalingEvent.version = unmarshaledData.Version
	brokenMarshalingEvent.occurredAt = unmarshaledData.OccurredAt

	return nil
}

/*** mocked Event with broken unmarshaling ***/

type BrokenUnmarshalingEvent struct {
	id         *SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenUnmarshalingEvent(forId *SomeID, withVersion uint) *BrokenUnmarshalingEvent {
	return &BrokenUnmarshalingEvent{
		id:         forId,
		name:       "BrokenUnmarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) Identifier() string {
	return brokenUnmarshalingEvent.id.String()
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) EventName() string {
	return brokenUnmarshalingEvent.name
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) OccurredAt() string {
	return brokenUnmarshalingEvent.occurredAt
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) StreamVersion() uint {
	return brokenUnmarshalingEvent.version
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID         *SomeID `json:"CustomerID"`
		Name       string  `json:"name"`
		Version    uint    `json:"version"`
		OccurredAt string  `json:"occurredAt"`
	}{
		ID:         brokenUnmarshalingEvent.id,
		Name:       brokenUnmarshalingEvent.name,
		Version:    brokenUnmarshalingEvent.version,
		OccurredAt: brokenUnmarshalingEvent.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func (brokenUnmarshalingEvent *BrokenUnmarshalingEvent) UnmarshalJSON(data []byte) error {
	return errors.New("mocked marshaling error")
}

/*** Unmarshal mocked events ***/

func Unmarshal(name string, payload []byte, streamVersion uint) (shared.DomainEvent, error) {
	defaultErrFormat := "unmarshalDomainEvent [%s] failed: %w"

	switch name {
	case "SomeEvent":
		event := &SomeEvent{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	case "BrokenMarshalingEvent":
		event := &BrokenUnmarshalingEvent{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	case "BrokenUnmarshalingEvent":
		event := &BrokenUnmarshalingEvent{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	default:
		return nil, errors.New("unknown mocked event to unmarshal")
	}
}
