package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/mocks"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerClone(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)

		customer, err := domain.ReconstituteCustomerFrom(
			shared.DomainEvents{
				events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion),
			},
		)
		So(err, ShouldBeNil)

		Convey("When it is cloned", func() {
			clonedCustomer := customer.Clone()

			Convey("It should be equal to the original customer", func() {
				So(clonedCustomer, ShouldResemble, customer)
			})
		})
	})
}

func TestCustomerStreamVersion(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		newEmailAddress, err := values.NewEmailAddress("john+changed@doe.com")
		So(err, ShouldBeNil)
		newConfirmableEmailAddress := newEmailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(2)

		customer, err := domain.ReconstituteCustomerFrom(
			shared.DomainEvents{
				events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion),
				events.EmailAddressWasChanged(id, newConfirmableEmailAddress, currentStreamVersion),
			},
		)
		So(err, ShouldBeNil)

		Convey("When it's streamVersion is retrieved", func() {
			streamVersion := customer.StreamVersion()

			Convey("It should expose the expected version", func() {
				So(streamVersion, ShouldResemble, currentStreamVersion)
			})
		})
	})
}

func TestCustomerApply(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		newEmailAddress, err := values.NewEmailAddress("john+changed@doe.com")
		So(err, ShouldBeNil)
		newConfirmableEmailAddress := newEmailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)

		customer, err := domain.ReconstituteCustomerFrom(
			shared.DomainEvents{
				events.ItWasRegistered(id, confirmableEmailAddress, personName, 1),
			},
		)
		So(err, ShouldBeNil)

		Convey("When latestEvents are applied", func() {
			currentStreamVersion++

			customer.Apply(
				shared.DomainEvents{
					events.EmailAddressWasChanged(id, newConfirmableEmailAddress, currentStreamVersion),
				},
			)

			Convey("It should be in the expected state", func() {
				So(customer.StreamVersion(), ShouldResemble, currentStreamVersion)

				changeEmailAddress, err := commands.NewChangeEmailAddress(
					id.String(),
					newEmailAddress.EmailAddress(),
				)
				So(err, ShouldBeNil)

				err = customer.Execute(changeEmailAddress)
				So(err, ShouldBeNil)

				So(customer.RecordedEvents(false), ShouldHaveLength, 0)
			})
		})
	})
}

func TestReconstituteCustomerFromWithInvalidEventStream(t *testing.T) {
	Convey("When a Customer is reconstituted from an empty EventStream", t, func() {
		var eventStream shared.DomainEvents

		_, err := domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})

	Convey("When a Customer is reconstituted from an EventStream without a Registered event", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)

		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		eventStream := shared.DomainEvents{
			events.EmailAddressWasConfirmed(id, emailAddress, uint(2)),
		}

		_, err = domain.ReconstituteCustomerFrom(eventStream)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInvalidEventStream), ShouldBeTrue)
		})
	})
}

func TestCustomerExecuteInvalidCommand(t *testing.T) {
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
			err := customer.Execute(nilInterfaceCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When a nil pointer command is handled", func() {
			var nilCommand *commands.ConfirmEmailAddress
			err := customer.Execute(nilCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an empty command is handled", func() {
			emptyCommand := &commands.ConfirmEmailAddress{}
			err := customer.Execute(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an unknown command is handled", func() {
			unknownCommand := new(mocks.Command)
			unknownCommand.On("AggregateID").Return(values.GenerateID())
			unknownCommand.On("CommandName").Return("unknown")
			err := customer.Execute(unknownCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandCanNotBeHandled), ShouldBeTrue)
			})
		})
	})
}
