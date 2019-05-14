package events_test

import (
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestItWasRegistered(t *testing.T) {
	Convey("Given valid parameters as input", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		Convey("When a new Registered event is created", func() {
			registered := events.ItWasRegistered(id, confirmableEmailAddress, personName)

			Convey("It should succeed", func() {
				So(registered, ShouldNotBeNil)
				So(registered, ShouldHaveSameTypeAs, (*events.Registered)(nil))
			})
		})
	})
}

func TestRegisteredExposesExpectedValues(t *testing.T) {
	Convey("Given a Registered event", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		beforeItOccurred := time.Now()
		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName)
		afterItOccurred := time.Now()
		So(registered, ShouldNotBeNil)
		So(registered, ShouldHaveSameTypeAs, (*events.Registered)(nil))

		Convey("It should expose the expected values", func() {
			So(registered.ID(), ShouldResemble, id)
			So(registered.ConfirmableEmailAddress(), ShouldResemble, confirmableEmailAddress)
			So(registered.PersonName(), ShouldResemble, personName)
			So(registered.Identifier(), ShouldEqual, id.String())
			So(registered.EventName(), ShouldEqual, "CustomerRegistered")
			itOccurred, err := time.Parse(shared.DomainEventMetaTimestampFormat, registered.OccurredAt())
			So(err, ShouldBeNil)
			So(beforeItOccurred, ShouldHappenBefore, itOccurred)
			So(afterItOccurred, ShouldHappenAfter, itOccurred)
		})
	})
}

func TestRegisteredMarshalJSON(t *testing.T) {
	Convey("Given a Registered event", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName)
		So(registered, ShouldNotBeNil)
		So(registered, ShouldHaveSameTypeAs, (*events.Registered)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := registered.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})
		})
	})
}

func TestRegisteredUnmarshalJSON(t *testing.T) {
	Convey("Given a Registered event marshaled to json", t, func() {
		id := values.GenerateID()

		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName)
		So(registered, ShouldNotBeNil)
		So(registered, ShouldHaveSameTypeAs, (*events.Registered)(nil))

		data, err := registered.MarshalJSON()

		Convey("And when it is unmarshaled", func() {
			unmarshaled := &events.Registered{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original Registered event", func() {
				So(err, ShouldBeNil)
				So(registered, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to Registered event", func() {
			unmarshaled := &events.Registered{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
