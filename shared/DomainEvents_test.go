package shared_test

import (
	"go-iddd/shared"
	"go-iddd/shared/mocks"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDomainEventsFirstEventShouldBeOfSameTypeAs(t *testing.T) {
	Convey("Given not empty DomainEvents", t, func() {
		domainEvents := shared.DomainEvents{}
		expectedType := new(mocks.SomethingHappend)

		Convey("And given the first event is of expected type", func() {
			domainEvents = append(domainEvents, expectedType)

			Convey("It should succeed", func() {
				err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

				So(err, ShouldBeNil)
			})
		})

		Convey("And given the first event is not of expected type", func() {
			domainEvents = append(domainEvents, new(mocks.SomethingElseHappend))

			Convey("It should fail", func() {
				err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given empty DomainEvents", t, func() {
		domainEvents := shared.DomainEvents{}
		expectedType := new(mocks.SomethingHappend)

		Convey("It should fail", func() {
			err := domainEvents.FirstEventShouldBeOfSameTypeAs(expectedType)

			So(err, ShouldBeError)
		})
	})
}
