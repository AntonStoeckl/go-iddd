package readmodel_test

import (
	"database/sql"
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/readmodel/domain/customer"
	eventsReadModel "go-iddd/service/customer/application/readmodel/domain/customer/events"
	valuesReadModel "go-iddd/service/customer/application/readmodel/domain/customer/values"
	eventsWriteModel "go-iddd/service/customer/application/writemodel/domain/customer/events"
	valuesWriteModel "go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/eventstore"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres/database"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerQueryHandlerScenarios(t *testing.T) {
	diContainer := setUpDiContainerForCustomerQueryHandlerScenarios()
	customerEventStoreForWriteModel := diContainer.GetCustomerEventStoreForWriteModel()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerView customer.View

		customerIDWriteModel := valuesWriteModel.GenerateCustomerID()
		customerIDReadModel := valuesReadModel.RebuildCustomerID(customerIDWriteModel.ID())
		theEmailAddress := "fiona@gallagher.net"
		emailAddress := valuesWriteModel.RebuildEmailAddress(theEmailAddress)
		confirmationHash := valuesWriteModel.GenerateConfirmationHash(emailAddress.EmailAddress())
		theGivenName := "Fiona"
		theFamilyName := "Galagher"
		personName := valuesWriteModel.RebuildPersonName(theGivenName, theFamilyName)
		theNewEmailAddress := "fiona@pratt.net"
		newEmailAddress := valuesWriteModel.RebuildEmailAddress(theNewEmailAddress)
		newConfirmationHash := valuesWriteModel.GenerateConfirmationHash(newEmailAddress.EmailAddress())

		expectedCustomerView := customer.View{
			ID:                      customerIDReadModel.ID(),
			EmailAddress:            emailAddress.EmailAddress(),
			IsEmailAddressConfirmed: false,
			GivenName:               personName.GivenName(),
			FamilyName:              personName.FamilyName(),
			Version:                 1,
		}

		Convey("\nSCENARIO 1: Retrieving a registered Customer", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", theGivenName, theFamilyName, theEmailAddress), func() {
				GivenCustomerRegistered(customerIDWriteModel, emailAddress, confirmationHash, personName, customerEventStoreForWriteModel)

				Convey("When the Customer is retrieved by ID", func() {
					customerView, err = queryHandler.CustomerViewByID(customerIDReadModel)
					So(err, ShouldBeNil)

					Convey("Then the Customer view should be as expected", func() {
						So(customerView, ShouldResemble, expectedCustomerView)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Retrieving a registered Customer with a confirmed emailAddress", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", theGivenName, theFamilyName, theEmailAddress), func() {
				GivenCustomerRegistered(customerIDWriteModel, emailAddress, confirmationHash, personName, customerEventStoreForWriteModel)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(customerIDWriteModel, emailAddress, 2, customerEventStoreForWriteModel)

					Convey("When the Customer is retrieved by ID", func() {
						customerView, err = queryHandler.CustomerViewByID(customerIDReadModel)
						So(err, ShouldBeNil)

						Convey("Then the Customer view should be as expected", func() {
							expectedCustomerView.IsEmailAddressConfirmed = true
							expectedCustomerView.Version = 2

							So(customerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Retrieving a registered Customer with an emailAddress that was confirmed and then changed", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", theGivenName, theFamilyName, theEmailAddress), func() {
				GivenCustomerRegistered(customerIDWriteModel, emailAddress, confirmationHash, personName, customerEventStoreForWriteModel)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(customerIDWriteModel, emailAddress, 2, customerEventStoreForWriteModel)

					Convey("And given she changed her emailAddress", func() {
						GivenEmailAddressChanged(customerIDWriteModel, newEmailAddress, newConfirmationHash, 3, customerEventStoreForWriteModel)

						Convey("When the Customer is retrieved by ID", func() {
							customerView, err = queryHandler.CustomerViewByID(customerIDReadModel)
							So(err, ShouldBeNil)

							Convey("Then the Customer view should be as expected", func() {
								expectedCustomerView.EmailAddress = newEmailAddress.EmailAddress()
								expectedCustomerView.Version = 3

								So(customerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Trying to retrieve a Customer which does not exist", func() {
			Convey("When the Customer is retrieved by ID", func() {
				customerView, err := queryHandler.CustomerViewByID(valuesReadModel.GenerateCustomerID())

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(customerView, ShouldBeZeroValue)
				})
			})
		})

		Reset(func() {
			err := customerEventStoreForWriteModel.Delete(customerIDWriteModel)
			So(err, ShouldBeNil)
		})
	})
}

func GivenCustomerRegistered(
	id valuesWriteModel.CustomerID,
	emailAddress valuesWriteModel.EmailAddress,
	hash valuesWriteModel.ConfirmationHash,
	name valuesWriteModel.PersonName,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		eventsWriteModel.CustomerWasRegistered(
			id,
			emailAddress,
			hash,
			name,
			1,
		),
	}

	err := customerEventStore.CreateStreamFrom(recordedEvents, id)
	So(err, ShouldBeNil)
}

func GivenEmailAddressConfirmed(
	id valuesWriteModel.CustomerID,
	emailAddress valuesWriteModel.EmailAddress,
	streamVersion uint,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		eventsWriteModel.CustomerEmailAddressWasConfirmed(
			id,
			emailAddress,
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, id)
	So(err, ShouldBeNil)
}

func GivenEmailAddressChanged(
	id valuesWriteModel.CustomerID,
	emailAddress valuesWriteModel.EmailAddress,
	confirmationHash valuesWriteModel.ConfirmationHash,
	streamVersion uint,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		eventsWriteModel.CustomerEmailAddressWasChanged(
			id,
			emailAddress,
			confirmationHash,
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, id)
	So(err, ShouldBeNil)
}

func setUpDiContainerForCustomerQueryHandlerScenarios() *cmd.DIContainer {
	config, err := cmd.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		panic(err)
	}

	err = migrator.Up()
	if err != nil {
		panic(err)
	}

	diContainer, err := cmd.NewDIContainer(
		db,
		eventsWriteModel.UnmarshalCustomerEvent,
		eventsReadModel.UnmarshalCustomerEvent,
	)

	if err != nil {
		panic(err)
	}

	return diContainer
}
