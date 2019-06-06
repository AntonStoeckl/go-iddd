package shared

import (
	"sort"
	"sync"

	"golang.org/x/xerrors"
)

type inMemoryEventStore struct {
	streamName string
	events     map[string]map[string]map[uint]DomainEvent
	eventsMux  sync.Mutex
}

func NewInMemoryEventStore(streamName string) *inMemoryEventStore {
	return &inMemoryEventStore{
		streamName: streamName,
		events:     make(map[string]map[string]map[uint]DomainEvent),
	}
}

/***** Implement shared.EventStore *****/

func (store *inMemoryEventStore) AppendToStream(
	identifier AggregateID,
	events DomainEvents,
) error {

	store.eventsMux.Lock()
	defer store.eventsMux.Unlock()

	// fist pass - assert that we have no concurrency conflict
	for _, event := range events {
		id := identifier.String()
		version := event.StreamVersion()

		if _, found := store.events[store.streamName][id][version]; found {
			return xerrors.Errorf("inMemoryEventStore.AppendToStream: %w", ErrConcurrencyConflict)
		}
	}

	// second pass - actually store the events
	if store.events[store.streamName] == nil {
		store.events[store.streamName] = make(map[string]map[uint]DomainEvent)
	}

	if store.events[store.streamName][identifier.String()] == nil {
		store.events[store.streamName][identifier.String()] = make(map[uint]DomainEvent)
	}

	for _, event := range events {
		store.events[store.streamName][identifier.String()][event.StreamVersion()] = event
	}

	return nil
}

func (store *inMemoryEventStore) LoadEventStream(identifier AggregateID) (DomainEvents, error) {
	store.eventsMux.Lock()
	defer store.eventsMux.Unlock()

	var eventStream DomainEvents

	events, found := store.events[store.streamName][identifier.String()]
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
	identifier AggregateID,
	fromVersion uint,
	maxEvents uint,
) (DomainEvents, error) {

	var eventStream DomainEvents
	var numEvents uint

	// err ignored: because this inMemory implementation of EventStore.LoadEventStream() can't return any error
	events, _ := store.LoadEventStream(identifier)

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
