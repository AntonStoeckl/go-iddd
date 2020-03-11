package eventstore

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"math"

	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type CustomerEventStore struct {
	eventStore es.EventStore
}

func NewCustomerEventStore(eventStore es.EventStore) *CustomerEventStore {
	return &CustomerEventStore{eventStore: eventStore}
}

func (customer *CustomerEventStore) EventStreamFor(id customer.ID) (es.DomainEvents, error) {
	eventStream, err := customer.eventStore.LoadEventStream(customer.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, "customerEventStore.EventStreamFor")
	}

	if len(eventStream) == 0 {
		err := errors.New("found empty event stream")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "customerEventStore.EventStreamFor")
	}

	return eventStream, nil
}

func (customer *CustomerEventStore) streamID(id customer.ID) es.StreamID {
	return es.NewStreamID(streamPrefix + "-" + id.ID())
}
