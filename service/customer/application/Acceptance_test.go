package application_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/mocked"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/query"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type acceptanceTestArtifacts struct {
	emailAddress    string
	givenName       string
	familyName      string
	newEmailAddress string
	newGivenName    string
	newFamilyName   string
}

func TestCustomerAcceptanceScenarios_ForRegisteringCustomers(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	eventStore := diContainer.GetCustomerEventStore()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var otherCustomerID values.CustomerID
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			emailAddress:    "fiona@gallagher.net",
			givenName:       "Fiona",
			familyName:      "Gallagher",
			newEmailAddress: "fiona@lishman.net",
		}

		Convey("\nSCENARIO: A prospective Customer registers her account", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, err = commandHandler.RegisterCustomer(aa.emailAddress, aa.givenName, aa.familyName)
				So(err, ShouldBeNil)

				expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
				details := fmt.Sprintf("\n\tGivenName: %s", expectedCustomerView.GivenName)
				details += fmt.Sprintf("\n\tFamilyName: %s", expectedCustomerView.FamilyName)
				details += fmt.Sprintf("\n\tEmailAddress: %s", expectedCustomerView.EmailAddress)
				details += fmt.Sprintf("\n\tIsEmailAddressConfirmed: %t", expectedCustomerView.IsEmailAddressConfirmed)

				Convey(fmt.Sprintf("Then her account should show the data she supplied: %s", details), func() {
					actualCustomerView = retrieveAccountData(queryHandler, customerID)
					So(actualCustomerView, ShouldResemble, expectedCustomerView)
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can't register because email address is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
					_, err = commandHandler.RegisterCustomer(aa.emailAddress, aa.givenName, aa.familyName)

					Convey(fmt.Sprintf("Then she should receive an error"), func() {
						So(err, ShouldBeError)
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can register with an email address that is not used any more", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey("And given the first Customer deleted her account", func() {
					err = commandHandler.DeleteCustomer(customerID.ID())
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						otherCustomerID, err = commandHandler.RegisterCustomer(aa.emailAddress, aa.givenName, aa.familyName)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})

				Convey("And given the first Customer changed her email address", func() {
					err = commandHandler.ChangeCustomerEmailAddress(customerID.ID(), aa.newEmailAddress)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						otherCustomerID, err = commandHandler.RegisterCustomer(aa.emailAddress, aa.givenName, aa.familyName)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)

			err = diContainer.GetCustomerEventStore().Delete(otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForConfirmingCustomerEmailAddresses(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	eventStore := diContainer.GetCustomerEventStore()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var confirmationHash values.ConfirmationHash
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			emailAddress:    "kevin@ball.net",
			givenName:       "Kevin",
			familyName:      "Ball",
			newEmailAddress: "levinia@ball.net",
		}

		Convey("\nSCENARIO: A Customer confirms his email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, confirmationHash = givenCustomerRegistered(aa, eventStore)

				Convey("When he confirms his email address", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), confirmationHash.Hash())
					So(err, ShouldBeNil)

					Convey("Then his email address should be confirmed", func() {
						actualCustomerView = retrieveAccountData(queryHandler, customerID)
						expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
						expectedCustomerView.IsEmailAddressConfirmed = true
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey("And when he confirms his email address again", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), confirmationHash.Hash())
							So(err, ShouldBeNil)

							Convey("Then his email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't confirm his email address, because the confirmation hash is invalid", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey("When he tries to confirm his email address with a wrong confirmation hash", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), "invalid_confirmation_hash")

					Convey("Then he should receive an error", func() {
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("And his email address should still be unconfirmed", func() {
							actualCustomerView = retrieveAccountData(queryHandler, customerID)
							expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer fails to confirm his already confirmed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey("And given he confirmed his email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2, eventStore)

					Convey("When he tries to confirm his email address again with a wrong confirmation hash", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), "invalid_confirmation_hash")

						Convey("Then he should receive an error", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("And his email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
								expectedCustomerView.IsEmailAddressConfirmed = true
								expectedCustomerView.Version = 3
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms his changed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey("And given he confirmed his email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2, eventStore)

					Convey(fmt.Sprintf("And given he changed his email address to [%s]", aa.newEmailAddress), func() {
						confirmationHash = givenCustomerEmailAddressWasChanged(customerID, aa, 3, eventStore)

						Convey("When he confirms his changed email address", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), confirmationHash.Hash())
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then his email address should be [%s] and confirmed", aa.newEmailAddress), func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
								expectedCustomerView.EmailAddress = aa.newEmailAddress
								expectedCustomerView.IsEmailAddressConfirmed = true
								expectedCustomerView.Version = 4
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerEmailAddresses(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	eventStore := diContainer.GetCustomerEventStore()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var otherCustomerID values.CustomerID
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			emailAddress:    "veronica@fisher.net",
			givenName:       "Veronica",
			familyName:      "Fisher",
			newEmailAddress: "veronica@pratt.net",
		}

		Convey("\nSCENARIO: A Customer changes her (confirmed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2, eventStore)

					Convey(fmt.Sprintf("When she changes her email address to [%s]", aa.newEmailAddress), func() {
						err = commandHandler.ChangeCustomerEmailAddress(customerID.ID(), aa.newEmailAddress)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("Then her email address should be [%s] and unconfirmed", aa.newEmailAddress), func() {
							actualCustomerView = retrieveAccountData(queryHandler, customerID)
							expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
							expectedCustomerView.EmailAddress = aa.newEmailAddress
							expectedCustomerView.Version = 3
							So(actualCustomerView, ShouldResemble, expectedCustomerView)

							Convey(fmt.Sprintf("And when she tries to change her email address to [%s] again", aa.newEmailAddress), func() {
								err = commandHandler.ChangeCustomerEmailAddress(customerID.ID(), aa.newEmailAddress)
								So(err, ShouldBeNil)

								Convey(fmt.Sprintf("Then her email address should still be [%s]", aa.newEmailAddress), func() {
									actualCustomerView = retrieveAccountData(queryHandler, customerID)
									So(actualCustomerView, ShouldResemble, expectedCustomerView)
								})
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't change her email address because it is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey(fmt.Sprintf("And given she changed here email address to [%s]", aa.newEmailAddress), func() {
					_ = givenCustomerEmailAddressWasChanged(customerID, aa, 2, eventStore)

					Convey(fmt.Sprintf("And given another Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
						otherCustomerID, _ = givenCustomerRegistered(aa, eventStore)

						Convey(fmt.Sprintf("When she also tries to change her email address to [%s]", aa.newEmailAddress), func() {
							err = commandHandler.ChangeCustomerEmailAddress(otherCustomerID.ID(), aa.newEmailAddress)

							Convey("Then she should receive an error", func() {
								So(err, ShouldBeError)
								So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
							})
						})
					})
				})
			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)

			err = diContainer.GetCustomerEventStore().Delete(otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerNames(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	eventStore := diContainer.GetCustomerEventStore()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			emailAddress:  "mikhailo@milkovich.net",
			givenName:     "Mikhailo",
			familyName:    "Milkovich",
			newGivenName:  "Mickey",
			newFamilyName: "Milkovich",
		}

		Convey("\nSCENARIO: A Customer changes his name", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa, eventStore)

				Convey(fmt.Sprintf("When he changes his name to [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
					err = commandHandler.ChangeCustomerName(customerID.ID(), aa.newGivenName, aa.newFamilyName)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then his name should be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
						actualCustomerView = retrieveAccountData(queryHandler, customerID)
						expectedCustomerView = buildExpectedCustomerViewForAcceptanceTest(customerID, aa)
						expectedCustomerView.GivenName = aa.newGivenName
						expectedCustomerView.FamilyName = aa.newFamilyName
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey(fmt.Sprintf("And when he tries to change his name to [%s %s] again", aa.newGivenName, aa.newFamilyName), func() {
							err = commandHandler.ChangeCustomerName(customerID.ID(), aa.newGivenName, aa.newFamilyName)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then his name should still be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
								actualCustomerView = retrieveAccountData(queryHandler, customerID)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForDeletingCustomers(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()
	eventStore := diContainer.GetCustomerEventStore()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var confirmationHash values.ConfirmationHash
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			emailAddress:    "karen@jackson.net",
			givenName:       "Karen",
			familyName:      "Jackson",
			newEmailAddress: "fiona@silverman.net",
			newGivenName:    "Karen",
			newFamilyName:   "Silverman",
		}

		Convey("\nSCENARIO: A Customer deletes her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, confirmationHash = givenCustomerRegistered(aa, eventStore)

				Convey("When she deletes her account", func() {
					err = commandHandler.DeleteCustomer(customerID.ID())
					So(err, ShouldBeNil)

					Convey("And when she tries to retrieve her account data", func() {
						actualCustomerView, err = queryHandler.CustomerViewByID(customerID.ID())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})

						Convey("And when she tries to delete her account again", func() {
							err = commandHandler.DeleteCustomer(customerID.ID())
							So(err, ShouldBeNil)

							Convey("Then her account should still be deleted", func() {
								actualCustomerView, err = queryHandler.CustomerViewByID(customerID.ID())
								So(err, ShouldBeError)
								So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
								So(actualCustomerView, ShouldBeZeroValue)
							})
						})
					})

					Convey("And when she tries to confirm her email address", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), confirmationHash.Hash())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her email address", func() {
						err = commandHandler.ChangeCustomerEmailAddress(customerID.ID(), aa.newEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her name", func() {
						err = commandHandler.ChangeCustomerName(customerID.ID(), aa.newGivenName, aa.newFamilyName)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})

			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_WhenCustomerWasNeverRegistered(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var actualCustomerView customer.View

		aa := acceptanceTestArtifacts{
			newEmailAddress: "fiona@pratt.net",
			newGivenName:    "Fiona",
			newFamilyName:   "Pratt",
		}

		customerID := values.GenerateCustomerID()
		confirmationHash := values.RebuildConfirmationHash(aa.newEmailAddress)

		Convey("\nSCENARIO: A hacker tries to get a non existing Customer account by guessing IDs", func() {
			Convey("When he tries to retrieve data for a non existing account", func() {
				actualCustomerView, err = queryHandler.CustomerViewByID(customerID.ID())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(actualCustomerView, ShouldBeZeroValue)
				})
			})

			Convey("And when he tries to confirm an email address", func() {
				err = commandHandler.ConfirmCustomerEmailAddress(customerID.ID(), confirmationHash.Hash())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to change an email address", func() {
				err = commandHandler.ChangeCustomerEmailAddress(customerID.ID(), aa.newEmailAddress)

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to change a name", func() {
				err = commandHandler.ChangeCustomerName(customerID.ID(), aa.newGivenName, aa.newFamilyName)

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to delete an account", func() {
				err = commandHandler.DeleteCustomer(customerID.ID())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}

func TestCustomerAcceptanceScenarios_InvalidInput(t *testing.T) {
	customerEventStoreMock := new(mocked.ForStoringCustomerEvents)
	commandHandlerWithMock := command.NewCustomerCommandHandler(customerEventStoreMock)

	Convey("Prepare test artifacts", t, func() {
		var err error

		aa := acceptanceTestArtifacts{
			emailAddress: "fiona@gallagher.net",
			givenName:    "Fiona",
			familyName:   "Gallagher",
		}

		customerID := values.GenerateCustomerID()
		confirmationHash := values.GenerateConfirmationHash(aa.emailAddress)
		invalidEmailAddress := "fiona@galagher.c"

		Convey("\nSCENARIO: A Customer supplies invalid input", func() {
			Convey(fmt.Sprintf("When she tries to register with an invalid email address [%s]", invalidEmailAddress), func() {
				_, err = commandHandlerWithMock.RegisterCustomer(invalidEmailAddress, aa.givenName, aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to register with an empty givenName", func() {
				_, err = commandHandlerWithMock.RegisterCustomer(aa.emailAddress, "", aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to register with an empty familyName", func() {
				_, err = commandHandlerWithMock.RegisterCustomer(aa.emailAddress, aa.givenName, "")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to confirm her email address with an empty id", func() {
				err = commandHandlerWithMock.ConfirmCustomerEmailAddress("", confirmationHash.Hash())

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to confirm her email address with an empty confirmation hash", func() {
				err = commandHandlerWithMock.ConfirmCustomerEmailAddress(customerID.ID(), "")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to change her email address with an empty id", func() {
				err = commandHandlerWithMock.ChangeCustomerEmailAddress("", aa.emailAddress)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey(fmt.Sprintf("When she tries to change her email address with an invalid email address [%s]", invalidEmailAddress), func() {
				err = commandHandlerWithMock.ChangeCustomerEmailAddress(customerID.ID(), invalidEmailAddress)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to change her name with an empty id", func() {
				err = commandHandlerWithMock.ChangeCustomerName("", aa.givenName, aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to change her name with an empty given name", func() {
				err = commandHandlerWithMock.ChangeCustomerName(customerID.ID(), "", aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to change her name with an empty family name", func() {
				err = commandHandlerWithMock.ChangeCustomerName(customerID.ID(), aa.givenName, "")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she tries to delete her account with an empty id", func() {
				err = commandHandlerWithMock.DeleteCustomer("")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})
		})
	})
}

func givenCustomerRegistered(
	aa acceptanceTestArtifacts,
	eventStore *eventstore.CustomerEventStore,
) (values.CustomerID, values.ConfirmationHash) {

	customerID := values.GenerateCustomerID()
	emailAddress := values.RebuildEmailAddress(aa.emailAddress)
	confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
	personName := values.RebuildPersonName(aa.givenName, aa.familyName)

	event := events.CustomerWasRegistered(
		customerID,
		emailAddress,
		confirmationHash,
		personName,
		1,
	)

	err := eventStore.CreateStreamFrom(es.DomainEvents{event}, customerID)
	So(err, ShouldBeNil)

	return customerID, confirmationHash
}

func givenCustomerEmailAddressWasConfirmed(
	customerID values.CustomerID,
	aa acceptanceTestArtifacts,
	streamVersion uint,
	eventStore *eventstore.CustomerEventStore,
) {

	emailAddress := values.RebuildEmailAddress(aa.emailAddress)

	event := events.CustomerEmailAddressWasConfirmed(
		customerID,
		emailAddress,
		streamVersion,
	)

	err := eventStore.Add(es.DomainEvents{event}, customerID)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	customerID values.CustomerID,
	aa acceptanceTestArtifacts,
	streamVersion uint,
	eventStore *eventstore.CustomerEventStore,
) values.ConfirmationHash {

	emailAddress := values.RebuildEmailAddress(aa.newEmailAddress)
	previousEmailAddress := values.RebuildEmailAddress(aa.emailAddress)
	confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())

	event := events.CustomerEmailAddressWasChanged(
		customerID,
		emailAddress,
		confirmationHash,
		previousEmailAddress,
		streamVersion,
	)

	err := eventStore.Add(es.DomainEvents{event}, customerID)
	So(err, ShouldBeNil)

	return confirmationHash
}

func retrieveAccountData(
	queryHandler *query.CustomerQueryHandler,
	id values.CustomerID,
) customer.View {

	customerView, err := queryHandler.CustomerViewByID(id.ID())
	So(err, ShouldBeNil)

	return customerView
}

func buildExpectedCustomerViewForAcceptanceTest(customerID values.CustomerID, aa acceptanceTestArtifacts) customer.View {
	return customer.View{
		ID:                      customerID.ID(),
		EmailAddress:            aa.emailAddress,
		IsEmailAddressConfirmed: false,
		GivenName:               aa.givenName,
		FamilyName:              aa.familyName,
		Version:                 1,
	}
}
