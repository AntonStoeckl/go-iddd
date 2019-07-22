package shared

import (
	"sort"
	"sync"

	"golang.org/x/xerrors"
)

type inMemoryEventStore struct {
	streamName   string
	events       map[string]map[string]map[uint]DomainEvent
	eventsMux    sync.Mutex
	failOnceWith error
}

func NewInMemoryEventStore(streamName string) *inMemoryEventStore {
	return &inMemoryEventStore{
		streamName: streamName,
		events:     make(map[string]map[string]map[uint]DomainEvent),
	}
}

/***** Implement shared.EventStore *****/

func (store *inMemoryEventStore) AppendToStream(streamID *StreamID, events DomainEvents) error {
	store.eventsMux.Lock()
	defer store.eventsMux.Unlock()

	if store.failOnceWith != nil {
		return store.failedOnceWith()
	}

	// fist pass - assert that we have no concurrency conflict
	for _, event := range events {
		if _, found := store.events[store.streamName][streamID.String()][event.StreamVersion()]; found {
			return xerrors.Errorf("inMemoryEventStore.AppendToStream: %w", ErrConcurrencyConflict)
		}
	}

	// second pass - actually store the events
	if store.events[store.streamName] == nil {
		store.events[store.streamName] = make(map[string]map[uint]DomainEvent)
	}

	for _, event := range events {
		store.ensureIndex(streamID.String())
		store.events[store.streamName][streamID.String()][event.StreamVersion()] = event
	}

	return nil
}

func (store *inMemoryEventStore) ensureIndex(id string) {
	if store.events[store.streamName][id] == nil {
		store.events[store.streamName][id] = make(map[uint]DomainEvent)
	}
}

func (store *inMemoryEventStore) LoadEventStream(streamID *StreamID) (DomainEvents, error) {
	store.eventsMux.Lock()
	defer store.eventsMux.Unlock()

	if store.failOnceWith != nil {
		return nil, store.failedOnceWith()
	}

	var eventStream DomainEvents

	events, found := store.events[store.streamName][streamID.String()]
	if !found {
		return eventStream, nil
	}

	versions := make([]int, 0, len(events))

	for _, event := range events {
		versions = append(versions, int(event.StreamVersion()))
	}

	sort.Ints(versions)

	for _, version := range versions {
		eventStream = append(eventStream, events[uint(version)])
	}

	return eventStream, nil
}

func (store *inMemoryEventStore) LoadPartialEventStream(
	streamID *StreamID,
	fromVersion uint,
	maxEvents uint,
) (DomainEvents, error) {

	var eventStream DomainEvents
	var numEvents uint

	events, err := store.LoadEventStream(streamID)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		// skip versions smaller than fromVersion
		if uint(event.StreamVersion()) < fromVersion {
			continue
		}

		// stop if it has reached maxEvents
		if numEvents == maxEvents {
			break
		}

		eventStream = append(eventStream, event)
		numEvents++
	}

	return eventStream, nil
}

/***** For mocking errors in tests *****/

func (store *inMemoryEventStore) FailOnceWith(err error) {
	store.failOnceWith = err
}

func (store *inMemoryEventStore) failedOnceWith() error {
	err := store.failOnceWith
	store.failOnceWith = nil

	return err
}
