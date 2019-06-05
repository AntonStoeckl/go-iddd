package shared_test

import (
	"fmt"
	"go-iddd/shared"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestInMemoryEventStore(t *testing.T) {
	Convey("Given an empty event stream", t, func() {
		es := shared.NewInMemoryEventStore("test")
		id := &someID{id: uuid.New().String()}
		var expectedEventStream shared.DomainEvents

		Convey("When events are appended", func() {
			event1 := createSomeEvent(id, 1)
			event2 := createSomeEvent(id, 2)
			err := es.AppendToStream(id, shared.DomainEvents{event1, event2})
			expectedEventStream = append(expectedEventStream, event1, event2)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And when the eventStream is loaded", func() {
					loadedEventStream, err := es.LoadEventStream(id)

					Convey("It should contain the expected events", func() {
						So(err, ShouldBeNil)
						So(loadedEventStream, ShouldResemble, expectedEventStream)
					})

					Convey("And when further events are appended", func() {
						event3 := createSomeEvent(id, 3)
						event4 := createSomeEvent(id, 4)
						err := es.AppendToStream(id, shared.DomainEvents{event3, event4})
						expectedEventStream = append(expectedEventStream, event3, event4)

						Convey("It should succeed", func() {
							So(err, ShouldBeNil)

							Convey("And when the eventStream is loaded", func() {
								loadedEventStream, err := es.LoadEventStream(id)

								Convey("It should contain the expected events", func() {
									So(err, ShouldBeNil)
									So(loadedEventStream, ShouldResemble, expectedEventStream)
								})
							})
						})
					})
				})
			})

			Convey("And when events with a conflicting version are appended", func() {
				err := es.AppendToStream(id, shared.DomainEvents{event2})

				Convey("It should not append them", func() {
					So(xerrors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)
					loadedEventStream, err := es.LoadEventStream(id)
					So(err, ShouldBeNil)
					So(loadedEventStream, ShouldResemble, expectedEventStream)
				})
			})
		})

		Convey("And when the empty eventStream is loaded", func() {
			loadedEventStream, err := es.LoadEventStream(id)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(loadedEventStream, ShouldHaveLength, 0)
			})
		})
	})
}

/*** Test Helpers ***/

type someID struct {
	id string
}

func (someID *someID) String() string {
	return someID.id
}

func (someID *someID) Equals(other shared.AggregateID) bool {
	return true // not needed in scope of this test
}

type someEvent struct {
	id      *someID
	name    string
	version uint
}

func createSomeEvent(forId *someID, withVersion uint) *someEvent {
	return &someEvent{id: forId, name: fmt.Sprintf("testEvent%d", withVersion), version: withVersion}
}

func (someEvent *someEvent) Identifier() string {
	return someEvent.id.String()
}

func (someEvent *someEvent) EventName() string {
	return someEvent.name
}

func (someEvent *someEvent) OccurredAt() string {
	return time.Now().Format(shared.DomainEventMetaTimestampFormat)
}

func (someEvent *someEvent) StreamVersion() uint {
	return someEvent.version
}
