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

func (someID SomeID) String() string {
	return someID.ID
}

func (someID SomeID) Equals(other shared.IdentifiesAggregates) bool {
	_ = other

	return true // not needed in scope of this test
}

/*** mocked Event that works ***/

type SomeEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateSomeEvent(forId SomeID, withVersion uint) SomeEvent {
	return SomeEvent{
		id:         forId,
		name:       "SomeEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (someEvent SomeEvent) EventName() string {
	return someEvent.name
}

func (someEvent SomeEvent) OccurredAt() string {
	return someEvent.occurredAt
}

func (someEvent SomeEvent) StreamVersion() uint {
	return someEvent.version
}

func (someEvent SomeEvent) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID         string `json:"customerID"`
		Name       string `json:"name"`
		Version    uint   `json:"version"`
		OccurredAt string `json:"occurredAt"`
	}{
		ID:         someEvent.id.ID,
		Name:       someEvent.name,
		Version:    someEvent.version,
		OccurredAt: someEvent.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalSomeEventFromJSON(data []byte) SomeEvent {
	someEvent := SomeEvent{
		id:         SomeID{ID: jsoniter.Get(data, "customerID").ToString()},
		name:       jsoniter.Get(data, "name").ToString(),
		version:    jsoniter.Get(data, "version").ToUint(),
		occurredAt: jsoniter.Get(data, "occurredAt").ToString(),
	}

	return someEvent
}

/*** mocked Event with broken marshaling ***/

type BrokenMarshalingEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenMarshalingEvent(forId SomeID, withVersion uint) BrokenMarshalingEvent {
	return BrokenMarshalingEvent{
		id:         forId,
		name:       "BrokenMarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (brokenMarshalingEvent BrokenMarshalingEvent) EventName() string {
	return brokenMarshalingEvent.name
}

func (brokenMarshalingEvent BrokenMarshalingEvent) OccurredAt() string {
	return brokenMarshalingEvent.occurredAt
}

func (brokenMarshalingEvent BrokenMarshalingEvent) StreamVersion() uint {
	return brokenMarshalingEvent.version
}

func (brokenMarshalingEvent BrokenMarshalingEvent) MarshalJSON() ([]byte, error) {
	return nil, errors.New("mocked marshaling error")
}

/*** mocked Event with broken unmarshaling ***/

type BrokenUnmarshalingEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenUnmarshalingEvent(forId SomeID, withVersion uint) BrokenUnmarshalingEvent {
	return BrokenUnmarshalingEvent{
		id:         forId,
		name:       "BrokenUnmarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (brokenUnmarshalingEvent BrokenUnmarshalingEvent) EventName() string {
	return brokenUnmarshalingEvent.name
}

func (brokenUnmarshalingEvent BrokenUnmarshalingEvent) OccurredAt() string {
	return brokenUnmarshalingEvent.occurredAt
}

func (brokenUnmarshalingEvent BrokenUnmarshalingEvent) StreamVersion() uint {
	return brokenUnmarshalingEvent.version
}

/*** Unmarshal mocked events ***/

func Unmarshal(name string, payload []byte, streamVersion uint) (shared.DomainEvent, error) {
	_ = streamVersion
	defaultErrFormat := "unmarshalDomainEvent [%s] failed: %w"

	switch name {
	case "SomeEvent":
		return UnmarshalSomeEventFromJSON(payload), nil
	case "BrokenMarshalingEvent":
		return BrokenMarshalingEvent{}, nil
	case "BrokenUnmarshalingEvent":
		return nil, errors.Errorf(defaultErrFormat, name, errors.New("mocked marshaling error"))
	default:
		return nil, errors.New("unknown mocked event to unmarshal")
	}
}
