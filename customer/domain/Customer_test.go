package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestReconstituteCustomerFromWithInvalidEventStream(t *testing.T) {
	Convey("When a Customer is reconstituted from an empty EventStream", t, func() {
		var eventStream shared.EventStream

		_, err := domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(xerrors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})

	Convey("When a Customer is reconstituted from an EventStream without a Registered event", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)

		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		eventStream := shared.EventStream{
			events.EmailAddressWasConfirmed(id, emailAddress, uint(2)),
		}

		_, err = domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(xerrors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})
}

func TestCustomerApplyInvalidCommand(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.NewCustomerWith(register)

		Convey("When a nil interface command is handled", func() {
			var nilInterfaceCommand shared.Command
			err := customer.Apply(nilInterfaceCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When a nil pointer command is handled", func() {
			var nilCommand *commands.ConfirmEmailAddress
			err := customer.Apply(nilCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an empty command is handled", func() {
			emptyCommand := &commands.ConfirmEmailAddress{}
			err := customer.Apply(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an unknown command is handled", func() {
			unknownCommand := &unknownCommand{}
			err := customer.Apply(unknownCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandCanNotBeHandled), ShouldBeTrue)
			})
		})
	})
}

/*** Test Helpers ***/

type unknownCommand struct{}

func (c *unknownCommand) AggregateID() shared.AggregateID {
	return values.GenerateID()
}

func (c *unknownCommand) CommandName() string {
	return "unknown"
}

func findCustomerEventIn(recordedEvents shared.EventStream, expectedEvent shared.DomainEvent) shared.DomainEvent {
	for _, event := range recordedEvents {
		if reflect.TypeOf(event) == reflect.TypeOf(expectedEvent) {
			return event
		}
	}

	return nil
}
