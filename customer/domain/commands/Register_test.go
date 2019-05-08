package commands_test

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegister(t *testing.T) {
	Convey("Given valid ID, EmailAddress and PersonName", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		Convey("When a new Register command is created", func() {
			register, err := commands.NewRegister(id, emailAddress, personName)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(register, ShouldHaveSameTypeAs, (*commands.Register)(nil))
			})
		})

		Convey("Given that ID is nil instead", func() {
			var id *values.ID

			conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
		})

		Convey("Given that EmailAddress is nil intead", func() {
			var emailAddress *values.EmailAddress

			conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
		})

		Convey("Given that PersonName is nil instead", func() {
			var personName *values.PersonName

			conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
		})
	})
}

func conveyNewRegisterWithInvalidInput(
	id *values.ID,
	emailAddress *values.EmailAddress,
	personName *values.PersonName,
) {

	Convey("When a new Register command is created", func() {
		register, err := commands.NewRegister(id, emailAddress, personName)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(register, ShouldBeNil)
		})
	})
}

func TestRegisterExposesExpectedValues(t *testing.T) {
	Convey("Given a Register command", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		register, err := commands.NewRegister(id, emailAddress, personName)
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			So(register.ID(), ShouldResemble, id)
			So(register.EmailAddress(), ShouldResemble, emailAddress)
			So(register.PersonName(), ShouldResemble, personName)
			So(register.CommandName(), ShouldEqual, "Register")
			So(register.AggregateIdentifier(), ShouldResemble, id)
		})
	})
}
