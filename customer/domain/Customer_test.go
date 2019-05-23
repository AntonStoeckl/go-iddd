package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

/*** Tests for Factory methods ***/

func TestRegisterCustomer(t *testing.T) {
	Convey("When a Customer is registered", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.Register(register)

		Convey("It should succeed", func() {
			So(customer, ShouldNotBeNil)
			So(customer, ShouldImplement, (*domain.Customer)(nil))

			Convey("And it should not apply further Register commands", func() {
				err = customer.Apply(register)
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandCanNotBeHandled), ShouldBeTrue)
			})
		})
	})
}

func TestCustomerExposesExpectedValues(t *testing.T) {
	Convey("Given a  Customer", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.Register(register)

		Convey("It should expose the expected AggregateIdentifier", func() {
			So(customer.AggregateIdentifier().String(), ShouldEqual, id)
		})

		Convey("It should expose the expected AggregateName", func() {
			So(customer.AggregateName(), ShouldEqual, "Customer")
		})
	})
}

/***** Test Customer business cases (other than Register) *****/

func TestConfirmEmailAddressOfCustomer(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.Register(register)

		Convey("When emailAddress is confirmed with not matching hash", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(id, emailAddress, "some_not_matching_hash")
			So(err, ShouldBeNil)

			err = customer.Apply(confirmEmailAddress)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
			})
		})
	})
}

/***** Test applying invalid commands *****/

func TestCustomerApplyInvalidCommand(t *testing.T) {
	Convey("Given a  Customer", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"
		givenName := "John"
		familyName := "Doe"

		register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
		So(err, ShouldBeNil)

		customer := domain.Register(register)

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

type unknownCommand struct{}

func (c *unknownCommand) AggregateIdentifier() shared.AggregateIdentifier {
	return values.GenerateID()
}

func (c *unknownCommand) CommandName() string {
	return "unknown"
}
