package commands_test

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegister(t *testing.T) {
	Convey("Given that ID, ConfirmableEmailAddress and Name are valid", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		name := valueobjects.NewName("Anton", "Stöckl")

		Convey("When NewRegister is invoked", func() {
			register, err := commands.NewRegister(id, emailAddress, name)

			Convey("Then it should create a Register command", func() {
				So(err, ShouldBeNil)
				So(register, ShouldImplement, (*commands.Register)(nil))
			})

			Convey("And then it should expose the expected CommandName, Identifier, ID, ConfirmableEmailAddress and Name ", func() {
				So(register.CommandName(), ShouldEqual, "Register")
				So(register.Identifier(), ShouldEqual, id.String())
				So(register.ID(), ShouldEqual, id)
				So(register.ConfirmableEmailAddress(), ShouldEqual, emailAddress)
				So(register.Name(), ShouldEqual, name)
			})
		})
	})

	Convey("Given that ID is nil", t, func() {
		var id valueobjects.ID
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		name := valueobjects.NewName("Anton", "Stöckl")

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})

	Convey("Given that ConfirmableEmailAddress is nil", t, func() {
		id := valueobjects.GenerateID()
		var emailAddress valueobjects.ConfirmableEmailAddress
		name := valueobjects.NewName("Anton", "Stöckl")

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})

	Convey("Given that Name is nil", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		var name valueobjects.Name

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})
}

func conveyNewRegisterWithInvalidInput(
	id valueobjects.ID,
	emailAddress valueobjects.ConfirmableEmailAddress,
	name valueobjects.Name,
) {

	Convey("When NewRegister is invoked", func() {
		register, err := commands.NewRegister(id, emailAddress, name)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(register, ShouldBeNil)
		})
	})
}
