package mongodb

import (
	"context"
	"math"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

const streamPrefix = "customer"

type forRetrievingEventStreams func(ctx context.Context, streamID es.StreamID, fromVersion uint, maxEvents uint) (es.EventStream, error)
type forAppendingEventsToStreams func(streamID es.StreamID, events []es.DomainEvent, tx mongo.SessionContext) error
type forPurgingEventStreams func(streamID es.StreamID, tx mongo.SessionContext) error
type forAssertingUniqueEmailAddresses func(recordedEvents []es.DomainEvent, tx mongo.SessionContext) error
type forPurgingUniqueEmailAddresses func(customerID value.CustomerID, tx mongo.SessionContext) error

type CustomerMongodbEventStore struct {
	db                       *mongo.Collection
	retrieveEventStream      forRetrievingEventStreams
	appendEventsToStream     forAppendingEventsToStreams
	purgeEventStream         forPurgingEventStreams
	assertUniqueEmailAddress forAssertingUniqueEmailAddresses
	purgeUniqueEmailAddress  forPurgingUniqueEmailAddresses
}

func NewCustomerMongodbEventStore(
	db *mongo.Collection,
	retrieveEventStream forRetrievingEventStreams,
	appendEventsToStream forAppendingEventsToStreams,
	purgeEventStream forPurgingEventStreams,
	assertUniqueEmailAddress forAssertingUniqueEmailAddresses,
	purgeUniqueEmailAddress forPurgingUniqueEmailAddresses,
) application.EventStoreInterface {
	return &CustomerMongodbEventStore{
		db:                       db,
		retrieveEventStream:      retrieveEventStream,
		appendEventsToStream:     appendEventsToStream,
		purgeEventStream:         purgeEventStream,
		assertUniqueEmailAddress: assertUniqueEmailAddress,
		purgeUniqueEmailAddress:  purgeUniqueEmailAddress,
	}
}

func (s *CustomerMongodbEventStore) RetrieveEventStream(id value.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "CustomerMongodbEventStore.RetrieveEventStream"
	ctx, cancel := context.WithCancel(context.Background())
	eventStream, err := s.retrieveEventStream(ctx, s.streamID(id), 0, math.MaxUint32)
	defer cancel()
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("customer not found")
		return nil, shared.MarkAndWrapError(err, shared.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (s *CustomerMongodbEventStore) StartEventStream(customerRegistered domain.CustomerRegistered) error {
	var err error
	wrapWithMsg := "CustomerMongodbEventStore.StartEventStream"
	recordedEvents := []es.DomainEvent{customerRegistered}
	ctx, cancel := context.WithCancel(context.Background())
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err = s.assertUniqueEmailAddress(recordedEvents, sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			return nil, errors.Wrap(err, wrapWithMsg)
		}

		streamID := s.streamID(customerRegistered.CustomerID())

		if err = s.appendEventsToStream(streamID, recordedEvents, sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			if errors.Is(err, shared.ErrConcurrencyConflict) {
				return nil, shared.MarkAndWrapError(errors.New("found duplicate customer"), shared.ErrDuplicate, wrapWithMsg)
			}

			return nil, errors.Wrap(err, wrapWithMsg)
		}

		return "ok", nil
	}

	defer cancel()
	_, er := s.transact(ctx, callback)
	if er != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerMongodbEventStore) AppendToEventStream(recordedEvents es.RecordedEvents, id value.CustomerID) error {
	var err error
	wrapWithMsg := "CustomerMongodbEventStore.AppendToEventStream"

	ctx, cancel := context.WithCancel(context.Background())
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err = s.assertUniqueEmailAddress(recordedEvents, sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			return nil, errors.Wrap(err, wrapWithMsg)
		}

		if err = s.appendEventsToStream(s.streamID(id), recordedEvents, sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			return nil, errors.Wrap(err, wrapWithMsg)
		}

		return "ok", nil
	}
	defer cancel()
	_, er := s.transact(ctx, callback)
	if er != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerMongodbEventStore) PurgeEventStream(id value.CustomerID) error {
	var err error
	wrapWithMsg := "CustomerMongodbEventStore.PurgeEventStream"
	ctx, cancel := context.WithCancel(context.Background())
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err = s.purgeUniqueEmailAddress(id, sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			return nil, errors.Wrap(err, wrapWithMsg)
		}

		if err = s.purgeEventStream(s.streamID(id), sessCtx); err != nil {
			sessCtx.AbortTransaction(sessCtx)

			return nil, errors.Wrap(err, wrapWithMsg)
		}
		return "ok", nil
	}
	defer cancel()
	_, er := s.transact(ctx, callback)
	if er != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerMongodbEventStore) streamID(id value.CustomerID) es.StreamID {
	return es.BuildStreamID(streamPrefix + "-" + id.String())
}

func (s *CustomerMongodbEventStore) transact(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := s.db.Database().Client().StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return nil, err
	}
	return session.WithTransaction(ctx, fn)
}
