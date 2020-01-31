package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_Register(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := infrastructure.SetUpDIContainer()
		So(err, ShouldBeNil)
		commandHandler := diContainer.GetCustomerCommandHandler()

		register, err := commands.NewRegister(
			"john@doe.com",
			"John",
			"Doe",
		)
		So(err, ShouldBeNil)

		Convey("Given no Customer with the same ID exists", func() {
			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
				})

				err := diContainer.GetCustomerEventStore().Delete(register.CustomerID())
				So(err, ShouldBeNil)
			})
		})

		Convey("Given a Customer with the same ID already exists", func() {
			err = commandHandler.Register(register)

			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
				})
			})

			err := diContainer.GetCustomerEventStore().Delete(register.CustomerID())
			So(err, ShouldBeNil)
		})

		Convey("When a Customer is registered with an invalid command", func() {
			err := commandHandler.Register(commands.Register{})

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})
	})
}

func Test_Register_WithTransactionErrors(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)
		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		commandHandler := application.NewCommandHandler(customerEventStore, db)

		register, err := commands.NewRegister("john@doe.com", "John", "Doe")
		So(err, ShouldBeNil)

		Convey("Assuming that beginning the transaction fails", func() {
			dbMock.ExpectBegin().WillReturnError(lib.ErrTechnical)

			Convey("When a Customer is registered", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})

		Convey("Assuming that committing the transaction fails", func() {
			dbMock.ExpectBegin()
			dbMock.ExpectCommit().WillReturnError(lib.ErrTechnical)

			customerEventStore.On("CreateStreamFrom", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

			Convey("When a Customer is registered", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}
