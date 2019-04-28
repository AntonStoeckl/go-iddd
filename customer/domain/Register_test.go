package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegister(t *testing.T) {
	Convey("Given that ID, ConfirmableEmailAddress and PersonName are valid", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		personName := valueobjects.NewPersonName("Anton", "Stöckl")

		Convey("When NewRegister is invoked", func() {
			register, err := domain.NewRegister(id, emailAddress, personName)

			Convey("Then it should create a Register command", func() {
				So(err, ShouldBeNil)
				So(register, ShouldImplement, (*domain.Register)(nil))
			})

			Convey("And then it should expose the expected CommandName, AggregateIdentifier, ID, ConfirmableEmailAddress and PersonName ", func() {
				So(register.CommandName(), ShouldEqual, "Register")
				So(register.AggregateIdentifier(), ShouldEqual, id)
				So(register.ID(), ShouldEqual, id)
				So(register.ConfirmableEmailAddress(), ShouldEqual, emailAddress)
				So(register.PersonName(), ShouldEqual, personName)
			})
		})
	})

	Convey("Given that ID is nil", t, func() {
		var id valueobjects.ID
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		personName := valueobjects.NewPersonName("Anton", "Stöckl")

		conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
	})

	Convey("Given that ConfirmableEmailAddress is nil", t, func() {
		id := valueobjects.GenerateID()
		var emailAddress valueobjects.ConfirmableEmailAddress
		personName := valueobjects.NewPersonName("Anton", "Stöckl")

		conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
	})

	Convey("Given that PersonName is nil", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		var personName valueobjects.PersonName

		conveyNewRegisterWithInvalidInput(id, emailAddress, personName)
	})
}

func conveyNewRegisterWithInvalidInput(
	id valueobjects.ID,
	emailAddress valueobjects.ConfirmableEmailAddress,
	personName valueobjects.PersonName,
) {

	Convey("When NewRegister is invoked", func() {
		register, err := domain.NewRegister(id, emailAddress, personName)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(register, ShouldBeNil)
		})
	})
}
