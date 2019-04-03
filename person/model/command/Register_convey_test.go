package command_test

import (
	"go-iddd/person/model/command"
	"go-iddd/person/model/vo"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegister(t *testing.T) {
	id := vo.NewID("12345")
	emailAddress := vo.NewEmailAddress("foo@bar.com")
	name := vo.NewName("Anton", "Stöckl")

	Convey("Given ID, EmailAddress and Name are valid", t, func() {
		Convey("When NewRegister is called", func() {
			register, err := command.NewRegister(id, emailAddress, name)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("And it should return a Register Command", func() {
				So(register, ShouldImplement, (*command.Register)(nil))
			})

			Convey("And it's CommandName should be as expected", func() {
				So(register.CommandName(), ShouldEqual, "Register")
			})

			Convey("And it's Identifier should be as expected", func() {
				So(register.Identifier(), ShouldEqual, id.ID())
			})
		})
	})
}

func TestNewRegisterWithInvalidInput(t *testing.T) {
	id := vo.NewID("12345")
	emailAddress := vo.NewEmailAddress("foo@bar.com")
	name := vo.NewName("Anton", "Stöckl")

	conveyNewRegisterWithInvalidInput := func(id vo.ID, emailAddress vo.EmailAddress, name vo.Name) {
		Convey("When NewRegister is called", func() {
			register, err := command.NewRegister(id, emailAddress, name)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError)
			})

			Convey("And it should return nil instead of a Register Command", func() {
				So(register, ShouldBeNil)
			})
		})
	}

	Convey("Given ID is nil", t, func() {
		var id vo.ID

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})

	Convey("Given EmailAddress is nil", t, func() {
		var emailAddress vo.EmailAddress

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})

	Convey("Given Name is nil", t, func() {
		var name vo.Name

		conveyNewRegisterWithInvalidInput(id, emailAddress, name)
	})
}
