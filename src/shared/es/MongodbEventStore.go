package es

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbEventStore struct {
	db                   *mongo.Collection
	marshalDomainEvent   MarshalDomainEvent
	unmarshalDomainEvent UnmarshalDomainEvent
}

func NewMongodbEventStore(
	db *mongo.Collection,
	marshalDomainEvent MarshalDomainEvent,
	unmarshalDomainEvent UnmarshalDomainEvent,
) *MongodbEventStore {

	return &MongodbEventStore{
		marshalDomainEvent:   marshalDomainEvent,
		unmarshalDomainEvent: unmarshalDomainEvent,
		db:                   db,
	}
}

type eventFromDB struct {
	Payload       string `json:"payload" bson:"payload"`
	EventName     string `json:"event_name" bson:"event_name"`
	StreamVersion uint   `json:"stream_version" bson:"stream_version"`
}

func (s *MongodbEventStore) RetrieveEventStream(
	ctx context.Context,
	streamID StreamID,
	fromVersion uint,
	maxEvents uint,
) (EventStream, error) {

	var err error
	wrapWithMsg := "retrieveEventStream"

	// queryTemplate := `SELECT event_name, payload, stream_version FROM %name%
	// 					WHERE stream_id = $1 AND stream_version >= $2
	// 					ORDER BY stream_version ASC
	// 					LIMIT $3`

	// query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)
	project := primitive.D{
		primitive.E{Key: "event_name", Value: 1},
		primitive.E{Key: "payload", Value: 1},
		primitive.E{Key: "stream_version", Value: 1},
	}
	query := primitive.D{
		primitive.E{Key: "stream_id", Value: streamID.String()},
		primitive.E{Key: "stream_version", Value: primitive.D{
			primitive.E{Key: "$gte", Value: fromVersion},
		}},
	}
	opt := options.Find().SetSort(primitive.D{
		primitive.E{Key: "stream_version", Value: 1},
	}).SetLimit(int64(maxEvents)).SetProjection(project)
	eventRows, err := s.db.Find(ctx, query, opt)
	if err != nil {
		return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	defer eventRows.Close(ctx)

	var eventStream EventStream
	var domainEvent DomainEvent
	for eventRows.Next(ctx) {
		mdl := new(eventFromDB)
		if e := eventRows.Decode(mdl); e != nil {
			return nil, shared.MarkAndWrapError(e, shared.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = s.unmarshalDomainEvent(mdl.EventName, []byte(mdl.Payload), mdl.StreamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrUnmarshalingFailed, wrapWithMsg)
		}

		eventStream = append(eventStream, domainEvent)
	}

	return eventStream, nil
}

func (s *MongodbEventStore) AppendEventsToStream(
	streamID StreamID,
	events []DomainEvent,
	tx mongo.SessionContext,
) error {

	var err error
	wrapWithMsg := "appendEventsToStream"

	// queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, occurred_at, payload)
	// 					VALUES ($1, $2, $3, $4, $5)`
	// query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	for _, event := range events {
		var eventJSON []byte

		eventJSON, err = s.marshalDomainEvent(event)
		if err != nil {
			return shared.MarkAndWrapError(err, shared.ErrMarshalingFailed, wrapWithMsg)
		}
		_, err = s.db.InsertOne(tx, primitive.D{
			primitive.E{Key: "stream_id", Value: streamID.String()},
			primitive.E{Key: "stream_version", Value: event.Meta().StreamVersion()},
			primitive.E{Key: "event_name", Value: event.Meta().EventName()},
			primitive.E{Key: "occurred_at", Value: event.Meta().OccurredAt()},
			primitive.E{Key: "payload", Value: eventJSON},
		})

		if err != nil {
			return errors.Wrap(err, wrapWithMsg)
		}
	}

	return nil
}

func (s *MongodbEventStore) PurgeEventStream(
	streamID StreamID,
	tx mongo.SessionContext,
) error {

	// queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	// query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)
	if _, err := s.db.DeleteMany(tx, primitive.D{
		primitive.E{Key: "stream_id", Value: streamID.String()},
	}); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, "purgeEventStream")
	}

	return nil
}
