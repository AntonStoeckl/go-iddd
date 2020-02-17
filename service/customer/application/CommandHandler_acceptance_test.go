package application_test

import (
	"database/sql"
	"fmt"
	"go-iddd/service/cmd"
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

func TestCommandHandlerScenarios(t *testing.T) {
	diContainer := setUpDiContainerForCommandHandlerScenarios()
	customerEventStore := diContainer.GetCustomerEventStore()
	commandHandler := diContainer.GetCustomerCommandHandler()

	Convey("Prepare test artifacts", t, func() {
		emailAddress := "fiona@gallagher.net"
		givenName := "Fiona"
		familyName := "Gallagher"
		invalidConfirmationHash := values.GenerateConfirmationHash("foo@bar.com").Hash()
		changedEmailAddress := "fiona@pratt.net"

		register, err := commands.NewRegister(
			emailAddress,
			givenName,
			familyName,
		)
		So(err, ShouldBeNil)

		customerID := register.CustomerID()
		confirmationHash := register.ConfirmationHash().Hash()

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			emailAddress,
			confirmationHash,
		)
		So(err, ShouldBeNil)

		confirmEmailAddressWithInvalidHash, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			emailAddress,
			invalidConfirmationHash,
		)
		So(err, ShouldBeNil)

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			customerID.ID(),
			changedEmailAddress,
		)
		So(err, ShouldBeNil)

		changedConfirmationHash := changeEmailAddress.ConfirmationHash().Hash()

		confirmChangedEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			changedEmailAddress,
			changedConfirmationHash,
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: A prospective Customer registers", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				err := commandHandler.Register(register)
				So(err, ShouldBeNil)

				Convey("Then she should have an unconfirmed account", func() {
					ThenEventStreamShouldBe(
						es.DomainEvents{
							events.Registered{},
						},
						customerEventStore,
						customerID,
					)
				})
			})
		})

		Convey("\nSCENARIO 2: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey(fmt.Sprintf("When she confirms her email address with confirmation hash [%s]", confirmationHash), func() {
						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then her email address should be confirmed", func() {
							ThenEventStreamShouldBe(
								es.DomainEvents{
									events.Registered{},
									events.EmailAddressConfirmed{},
								},
								customerEventStore,
								customerID,
							)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: A Customer fails to confirm her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey(fmt.Sprintf("When she tries to confirm her email address with wrong confirmation hash [%s]", invalidConfirmationHash), func() {
						err = commandHandler.ConfirmEmailAddress(confirmEmailAddressWithInvalidHash)
						So(err, ShouldBeError)

						Convey("Then it should fail", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("and her email address should be unconfirmed", func() {
								ThenEventStreamShouldBe(
									es.DomainEvents{
										events.Registered{},
										events.EmailAddressConfirmationFailed{},
									},
									customerEventStore,
									customerID,
								)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: A Customer confirms her email address again with the right confirmationHash", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey("and she confirmed her email address", func() {
						GivenEmailAddressConfirmed(confirmEmailAddress, customerEventStore, 2)

						Convey(fmt.Sprintf("When she tries to confirm it again with confirmation hash [%s]", confirmationHash), func() {
							err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then it should be ignored", func() {
								ThenEventStreamShouldBe(
									es.DomainEvents{
										events.Registered{},
										events.EmailAddressConfirmed{},
									},
									customerEventStore,
									customerID,
								)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: A Customer confirms her email address again, but with a wrong confirmation hash", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", confirmationHash), func() {
					Convey("and she confirmed her email address", func() {
						GivenEmailAddressConfirmed(confirmEmailAddress, customerEventStore, 2)

						Convey(fmt.Sprintf("When she tries to confirm it again with confirmation hash [%s]", confirmationHash), func() {
							err = commandHandler.ConfirmEmailAddress(confirmEmailAddressWithInvalidHash)
							So(err, ShouldBeError)

							Convey("Then it should fail", func() {
								So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

								ThenEventStreamShouldBe(
									es.DomainEvents{
										events.Registered{},
										events.EmailAddressConfirmed{},
										events.EmailAddressConfirmationFailed{},
									},
									customerEventStore,
									customerID,
								)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 6: A Customer changes her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("When she changes her email address to [%s]", changedEmailAddress), func() {
					err = commandHandler.ChangeEmailAddress(changeEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then her email address should be changed", func() {
						ThenEventStreamShouldBe(
							es.DomainEvents{
								events.Registered{},
								events.EmailAddressChanged{},
							},
							customerEventStore,
							customerID,
						)
					})
				})
			})
		})

		Convey("\nSCENARIO 7: A Customer changes her email address twice", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey(fmt.Sprintf("and she changed her email address to [%s]", changedEmailAddress), func() {
					GivenEmailAddressChanged(changeEmailAddress, customerEventStore, 2)

					Convey(fmt.Sprintf("When she tries to change it again to [%s]", changedEmailAddress), func() {
						err = commandHandler.ChangeEmailAddress(changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then it should be ignored", func() {
							ThenEventStreamShouldBe(
								es.DomainEvents{
									events.Registered{},
									events.EmailAddressChanged{},
								},
								customerEventStore,
								customerID,
							)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 8: A Customer confirms her changed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", givenName, familyName, emailAddress), func() {
				GivenCustomerRegistered(register, customerEventStore)

				Convey("and she confirmed her email address", func() {
					GivenEmailAddressConfirmed(confirmEmailAddress, customerEventStore, 2)

					Convey(fmt.Sprintf("and she changed her email address to [%s]", changedEmailAddress), func() {
						GivenEmailAddressChanged(changeEmailAddress, customerEventStore, 3)

						Convey(fmt.Sprintf("and she was issued a confirmation hash [%s]", changedConfirmationHash), func() {
							Convey(fmt.Sprintf("When she confirms her changed email address with confirmation hash [%s]", changedConfirmationHash), func() {
								err = commandHandler.ConfirmEmailAddress(confirmChangedEmailAddress)
								So(err, ShouldBeNil)

								Convey("Then her changed email address should be confirmed", func() {
									ThenEventStreamShouldBe(
										es.DomainEvents{
											events.Registered{},
											events.EmailAddressConfirmed{},
											events.EmailAddressChanged{},
											events.EmailAddressConfirmed{},
										},
										customerEventStore,
										customerID,
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

func GivenCustomerRegistered(register commands.Register, customerEventStore *eventstore.CustomerEventStore) {
	recordedEvents := es.DomainEvents{
		events.ItWasRegistered(
			register.CustomerID(),
			register.EmailAddress(),
			register.ConfirmationHash(),
			register.PersonName(),
			1,
		),
	}

	err := customerEventStore.CreateStreamFrom(recordedEvents, register.CustomerID())
	So(err, ShouldBeNil)
}

func GivenEmailAddressConfirmed(
	confirmEmailAddress commands.ConfirmEmailAddress,
	customerEventStore *eventstore.CustomerEventStore,
	streamVersion uint,
) {

	recordedEvents := es.DomainEvents{
		events.EmailAddressWasConfirmed(
			confirmEmailAddress.CustomerID(),
			confirmEmailAddress.EmailAddress(),
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, confirmEmailAddress.CustomerID())
	So(err, ShouldBeNil)
}

func GivenEmailAddressChanged(
	changeEmailAddress commands.ChangeEmailAddress,
	customerEventStore *eventstore.CustomerEventStore,
	streamVersion uint,
) {

	recordedEvents := es.DomainEvents{
		events.EmailAddressWasChanged(
			changeEmailAddress.CustomerID(),
			changeEmailAddress.EmailAddress(),
			changeEmailAddress.ConfirmationHash(),
			streamVersion,
		),
	}

	err := customerEventStore.Add(recordedEvents, changeEmailAddress.CustomerID())
	So(err, ShouldBeNil)
}

func ThenEventStreamShouldBe(
	domainEvents es.DomainEvents,
	customerEventStore *eventstore.CustomerEventStore,
	customerID values.CustomerID,
) {

	eventStream, err := customerEventStore.EventStreamFor(customerID)
	So(err, ShouldBeNil)

	So(eventStream, ShouldHaveLength, len(domainEvents))

	for idx, event := range domainEvents {
		So(eventStream[idx], ShouldHaveSameTypeAs, event)
	}
}

func setUpDiContainerForCommandHandlerScenarios() *cmd.DIContainer {
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
