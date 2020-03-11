package readmodel_test

import (
	"database/sql"
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/readmodel/domain/customer/queries"
	"go-iddd/service/customer/application/writemodel/domain/customer/events"
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
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
	customerEventStore := diContainer.GetCustomerEventStore()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		customerID := values.GenerateCustomerID()
		theEmailAddress := "fiona@gallagher.net"
		emailAddress := values.RebuildEmailAddress(theEmailAddress)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		theGivenName := "Fiona"
		theFamilyName := "Galagher"
		personName := values.RebuildPersonName(theGivenName, theFamilyName)
		theNewEmailAddress := "fiona@pratt.net"
		newEmailAddress := values.RebuildEmailAddress(theNewEmailAddress)
		newConfirmationHash := values.GenerateConfirmationHash(newEmailAddress.EmailAddress())

		customerByIDQuery := queries.BuildCustomerByID(customerID.ID())

		expectedCustomerView := customer.View{
			ID:                      customerID.ID(),
			EmailAddress:            emailAddress.EmailAddress(),
			IsEmailAddressConfirmed: false,
			GivenName:               personName.GivenName(),
			FamilyName:              personName.FamilyName(),
			Version:                 1,
		}

		Convey("\nSCENARIO 1: Retrieving a registered Customer", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", theGivenName, theFamilyName, theEmailAddress), func() {
				GivenCustomerRegistered(customerID, emailAddress, confirmationHash, personName, customerEventStore)

				Convey("When the Customer is retrieved by ID", func() {
					customerView, err := queryHandler.CustomerViewByID(customerByIDQuery)
					So(err, ShouldBeNil)

					Convey("Then the Customer view should be as expected", func() {
						So(customerView, ShouldResemble, expectedCustomerView)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Retrieving a registered Customer with a confirmed emailAddress", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", theGivenName, theFamilyName, theEmailAddress), func() {
				GivenCustomerRegistered(customerID, emailAddress, confirmationHash, personName, customerEventStore)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(customerID, emailAddress, 2, customerEventStore)

					Convey("When the Customer is retrieved by ID", func() {
						customerView, err := queryHandler.CustomerViewByID(customerByIDQuery)
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
				GivenCustomerRegistered(customerID, emailAddress, confirmationHash, personName, customerEventStore)

				Convey("And given she confirmed her emailAddress", func() {
					GivenEmailAddressConfirmed(customerID, emailAddress, 2, customerEventStore)

					Convey("And given she changed her emailAddress", func() {
						GivenEmailAddressChanged(customerID, newEmailAddress, newConfirmationHash, 3, customerEventStore)

						Convey("When the Customer is retrieved by ID", func() {
							customerView, err := queryHandler.CustomerViewByID(customerByIDQuery)
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
				customerByIDQuery = queries.BuildCustomerByID(values.GenerateCustomerID().ID())
				customerView, err := queryHandler.CustomerViewByID(customerByIDQuery)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(customerView, ShouldBeZeroValue)
				})
			})
		})
	})
}

func GivenCustomerRegistered(
	id values.CustomerID,
	emailAddress values.EmailAddress,
	hash values.ConfirmationHash,
	name values.PersonName,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		events.CustomerWasRegistered(
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
	id values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		events.CustomerEmailAddressWasConfirmed(
			id,
			emailAddress,
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, id)
	So(err, ShouldBeNil)
}

func GivenEmailAddressChanged(
	id values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	streamVersion uint,
	customerEventStore *eventstore.CustomerEventStore,
) {

	recordedEvents := es.DomainEvents{
		events.CustomerEmailAddressWasChanged(
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

	diContainer, err := cmd.NewDIContainer(db, events.UnmarshalCustomerEvent)
	if err != nil {
		panic(err)
	}

	return diContainer
}
