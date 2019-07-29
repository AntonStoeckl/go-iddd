package shared_test

import (
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDomainEventsFirstEventShouldBeOfSameTypeAs(t *testing.T) {
	Convey("Given not empty DomainEvents", t, func() {
		domainEvents := shared.DomainEvents{}
		expectedType := new(SomethingHappend)

		Convey("And given the first event is of expected type", func() {
			domainEvents = append(domainEvents, expectedType)

			Convey("It should succeed", func() {
				err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

				So(err, ShouldBeNil)
			})
		})

		Convey("And given the first event is not of expected type", func() {
			domainEvents = append(domainEvents, new(SomethingElseHappend))

			Convey("It should fail", func() {
				err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given empty DomainEvents", t, func() {
		domainEvents := shared.DomainEvents{}
		expectedType := new(SomethingHappend)

		Convey("It should fail", func() {
			err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

			So(err, ShouldBeError)
		})
	})
}

/*** test helpers ***/

type SomethingHappend struct{}

func (event *SomethingHappend) Identifier() string {
	return "something"
}

func (event *SomethingHappend) EventName() string {
	return "something"
}

func (event *SomethingHappend) OccurredAt() string {
	return "something"
}

func (event *SomethingHappend) StreamVersion() uint {
	return 1
}

type SomethingElseHappend struct{}

func (event *SomethingElseHappend) Identifier() string {
	return "something"
}

func (event *SomethingElseHappend) EventName() string {
	return "something"
}

func (event *SomethingElseHappend) OccurredAt() string {
	return "something"
}

func (event *SomethingElseHappend) StreamVersion() uint {
	return 1
}
