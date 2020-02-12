package application_test

import (
	"database/sql"
	"fmt"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/eventstore"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres/database"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	commandHandler                     *application.CommandHandler
	customerEventStore                 *eventstore.CustomerEventStore
	customerID                         values.CustomerID
	emailAddress                       string
	givenName                          string
	familyName                         string
	confirmationHash                   string
	invalidConfirmationHash            string
	changedEmailAddress                string
	changedConfirmationHash            string
	register                           commands.Register
	confirmEmailAddress                commands.ConfirmEmailAddress
	confirmEmailAddressWithInvalidHash commands.ConfirmEmailAddress
	changeEmailAddress                 commands.ChangeEmailAddress
	confirmChangedEmailAddress         commands.ConfirmEmailAddress
)

func TestCommandHandlerScenarios(t *testing.T) {
	Convey("Customer Lifecycle Scenarios", t, func() {
		var err error

		setUpForCommandHandlerScenarios()

		Convey("\nSCENARIO 1: A prospective Customer registers", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("Then she should have an unconfirmed account", func() {
					ThenEventStreamShouldBe(events.Registered{})
				})
			})
		})

		Convey("\nSCENARIO 2: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey(fmt.Sprintf("When she confirms her email address with confirmation hash [%s]", confirmationHash), func() {
						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then her email address should be confirmed", func() {
							ThenEventStreamShouldBe(
								events.Registered{},
								events.EmailAddressConfirmed{},
							)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: A Customer fails to confirm her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey(fmt.Sprintf("When she tries to confirm her email address with invalid confirmation hash [%s]", invalidConfirmationHash), func() {
						err = commandHandler.ConfirmEmailAddress(confirmEmailAddressWithInvalidHash)
						So(err, ShouldBeError)

						Convey("Then it should fail", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("and her email address should be unconfirmed", func() {
								ThenEventStreamShouldBe(
									events.Registered{},
									events.EmailAddressConfirmationFailed{},
								)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: A Customer confirms her email address twice", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey("and she confirmed her email address", func() {
						GivenEmailAddressConfirmed(2)

						Convey(fmt.Sprintf("When she tries to confirm it again with confirmation hash [%s]", confirmationHash), func() {
							err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then it should be ignored", func() {
								ThenEventStreamShouldBe(
									events.Registered{},
									events.EmailAddressConfirmed{},
								)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: A Customer changes her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey(fmt.Sprintf("When she changes her email address to [%s]", changedEmailAddress), func() {
					err = commandHandler.ChangeEmailAddress(changeEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then her email address should be changed", func() {
						ThenEventStreamShouldBe(
							events.Registered{},
							events.EmailAddressChanged{},
						)
					})
				})
			})
		})

		Convey("\nSCENARIO 6: A Customer changes her email address twice", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey(fmt.Sprintf("and she changed her email address to [%s]", changedEmailAddress), func() {
					GivenEmailAddressChanged(2)

					Convey(fmt.Sprintf("When she tries to change it again to [%s]", changedEmailAddress), func() {
						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then it should be ignored", func() {
							ThenEventStreamShouldBe(
								events.Registered{},
								events.EmailAddressChanged{},
							)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 7: A Customer confirms her changed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered()

				Convey("and she confirmed her email address", func() {
					GivenEmailAddressConfirmed(2)

					Convey(fmt.Sprintf("and she changed her email address to [%s]", changedEmailAddress), func() {
						GivenEmailAddressChanged(3)

						Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", changedConfirmationHash), func() {
							Convey(fmt.Sprintf("When she confirms her changed email address with confirmation hash [%s]", changedConfirmationHash), func() {
								err = commandHandler.ConfirmEmailAddress(confirmChangedEmailAddress)
								So(err, ShouldBeNil)

								Convey("Then her changed email address should be confirmed", func() {
									ThenEventStreamShouldBe(
										events.Registered{},
										events.EmailAddressConfirmed{},
										events.EmailAddressChanged{},
										events.EmailAddressConfirmed{},
									)
								})
							})
						})
					})
				})
			})
		})

		Reset(func() {
			err := customerEventStore.Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func GivenCustomerRegistered() {
	recordedEvents := es.DomainEvents{
		events.ItWasRegistered(
			customerID,
			values.RebuildEmailAddress(emailAddress),
			values.RebuildConfirmationHash(confirmationHash),
			values.RebuildPersonName(givenName, familyName),
			1,
		),
	}

	err := customerEventStore.CreateStreamFrom(recordedEvents, customerID)
	So(err, ShouldBeNil)
}

func GivenEmailAddressConfirmed(streamVersion uint) {
	recordedEvents := es.DomainEvents{
		events.EmailAddressWasConfirmed(
			customerID,
			values.RebuildEmailAddress(emailAddress),
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, customerID)
	So(err, ShouldBeNil)
}

func GivenEmailAddressChanged(streamVersion uint) {
	recordedEvents := es.DomainEvents{
		events.EmailAddressWasChanged(
			customerID,
			values.RebuildEmailAddress(changedEmailAddress),
			values.RebuildConfirmationHash(changedConfirmationHash),
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, customerID)
	So(err, ShouldBeNil)
}

func ThenEventStreamShouldBe(domainEvents ...es.DomainEvent) {
	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)

	So(eventStream, ShouldHaveLength, len(domainEvents))

	for idx, event := range domainEvents {
		So(eventStream[idx], ShouldHaveSameTypeAs, event)
	}
}

func setUpForCommandHandlerScenarios() {
	var err error

	diContainer := setUpDiContainerForCommandHandlerScenarios()
	customerEventStore = diContainer.GetCustomerEventStore()
	commandHandler = diContainer.GetCustomerCommandHandler()

	emailAddress = "fiona@gallagher.net"
	givenName = "Fiona"
	familyName = "Gallagher"

	register, err = commands.NewRegister(
		emailAddress,
		givenName,
		familyName,
	)
	So(err, ShouldBeNil)

	customerID = register.CustomerID()
	confirmationHash = register.ConfirmationHash().Hash()

	confirmEmailAddress, err = commands.NewConfirmEmailAddress(
		customerID.ID(),
		emailAddress,
		confirmationHash,
	)
	So(err, ShouldBeNil)

	invalidConfirmationHash = values.GenerateConfirmationHash(emailAddress).Hash()

	confirmEmailAddressWithInvalidHash, err = commands.NewConfirmEmailAddress(
		customerID.ID(),
		emailAddress,
		invalidConfirmationHash,
	)
	So(err, ShouldBeNil)

	changedEmailAddress = "fiona@pratt.net"

	changeEmailAddress, err = commands.NewChangeEmailAddress(
		customerID.ID(),
		changedEmailAddress,
	)
	So(err, ShouldBeNil)

	changedConfirmationHash = changeEmailAddress.ConfirmationHash().Hash()

	confirmChangedEmailAddress, err = commands.NewConfirmEmailAddress(
		customerID.ID(),
		changedEmailAddress,
		changedConfirmationHash,
	)
	So(err, ShouldBeNil)
}

func setUpDiContainerForCommandHandlerScenarios() *cmd.DIContainer {
	config, err := cmd.NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	err = db.Ping()
	So(err, ShouldBeNil)

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	So(err, ShouldBeNil)

	err = migrator.Up()
	So(err, ShouldBeNil)

	diContainer, err := cmd.NewDIContainer(
		db,
		events.UnmarshalCustomerEvent,
	)
	So(err, ShouldBeNil)

	return diContainer
}
