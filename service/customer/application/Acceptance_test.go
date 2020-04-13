package application_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var acceptanceTestCustomerEventStore command.ForStoringCustomerEvents

type acceptanceTestCollaborators struct {
	registerCustomer            application.ForRegisteringCustomers
	confirmCustomerEmailAddress application.ForConfirmingCustomerEmailAddresses
	changeCustomerEmailAddress  application.ForChangingCustomerEmailAddresses
	changeCustomerName          application.ForChangingCustomerNames
	deleteCustomer              application.ForDeletingCustomers
	customerViewByID            application.ForRetrievingCustomerViews
}

type acceptanceTestArtifacts struct {
	emailAddress    string
	givenName       string
	familyName      string
	newEmailAddress string
	newGivenName    string
	newFamilyName   string
}

func TestCustomerAcceptanceScenarios_ForRegisteringCustomers(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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
				customerID, err = ac.registerCustomer(aa.emailAddress, aa.givenName, aa.familyName)
				So(err, ShouldBeNil)

				expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
				details := fmt.Sprintf("\n\tGivenName: %s", expectedCustomerView.GivenName)
				details += fmt.Sprintf("\n\tFamilyName: %s", expectedCustomerView.FamilyName)
				details += fmt.Sprintf("\n\tEmailAddress: %s", expectedCustomerView.EmailAddress)
				details += fmt.Sprintf("\n\tIsEmailAddressConfirmed: %t", expectedCustomerView.IsEmailAddressConfirmed)

				Convey(fmt.Sprintf("Then her account should show the data she supplied: %s", details), func() {
					actualCustomerView, err = ac.customerViewByID(customerID.String())
					So(err, ShouldBeNil)
					So(actualCustomerView, ShouldResemble, expectedCustomerView)
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can't register because email address is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
					_, err = ac.registerCustomer(aa.emailAddress, aa.givenName, aa.familyName)

					Convey(fmt.Sprintf("Then she should receive an error"), func() {
						So(err, ShouldBeError)
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can register with an email address that is not used any more", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey("And given the first Customer deleted her account", func() {
					err = ac.deleteCustomer(customerID.String())
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						otherCustomerID, err = ac.registerCustomer(aa.emailAddress, aa.givenName, aa.familyName)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})

				Convey(fmt.Sprintf("Or given the first Customer changed her email address to [%s]", aa.newEmailAddress), func() {
					err = ac.changeCustomerEmailAddress(customerID.String(), aa.newEmailAddress)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						otherCustomerID, err = ac.registerCustomer(aa.emailAddress, aa.givenName, aa.familyName)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer tries to register with invalid input", func() {
			invalidEmailAddress := "fiona@galagher.c"

			Convey(fmt.Sprintf("When she supplies an invalid email address [%s]", invalidEmailAddress), func() {
				_, err = ac.registerCustomer(invalidEmailAddress, aa.givenName, aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she supplies with an empty givenName", func() {
				_, err = ac.registerCustomer(aa.emailAddress, "", aa.familyName)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she supplies with an empty familyName", func() {
				_, err = ac.registerCustomer(aa.emailAddress, aa.givenName, "")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)

			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForConfirmingCustomerEmailAddresses(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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
				customerID, confirmationHash = givenCustomerRegistered(aa)

				Convey("When he confirms his email address", func() {
					err = ac.confirmCustomerEmailAddress(customerID.String(), confirmationHash.String())
					So(err, ShouldBeNil)

					Convey("Then his email address should be confirmed", func() {
						actualCustomerView, err = ac.customerViewByID(customerID.String())
						So(err, ShouldBeNil)
						expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
						expectedCustomerView.IsEmailAddressConfirmed = true
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey("And when he confirms his email address again", func() {
							err = ac.confirmCustomerEmailAddress(customerID.String(), confirmationHash.String())
							So(err, ShouldBeNil)

							Convey("Then his email address should still be confirmed", func() {
								actualCustomerView, err = ac.customerViewByID(customerID.String())
								So(err, ShouldBeNil)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't confirm his email address, because the confirmation hash is not matching", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey("When he tries to confirm his email address with a wrong confirmation hash", func() {
					err = ac.confirmCustomerEmailAddress(customerID.String(), "invalid_confirmation_hash")

					Convey("Then he should receive an error", func() {
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("And his email address should still be unconfirmed", func() {
							actualCustomerView, err = ac.customerViewByID(customerID.String())
							So(err, ShouldBeNil)
							expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer fails to confirm his already confirmed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey("And given he confirmed his email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2)

					Convey("When he tries to confirm his email address again with a wrong confirmation hash", func() {
						err = ac.confirmCustomerEmailAddress(customerID.String(), "invalid_confirmation_hash")

						Convey("Then he should receive an error", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("And his email address should still be confirmed", func() {
								actualCustomerView, err = ac.customerViewByID(customerID.String())
								So(err, ShouldBeNil)
								expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
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
				customerID, _ = givenCustomerRegistered(aa)

				Convey("And given he confirmed his email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2)

					Convey(fmt.Sprintf("And given he changed his email address to [%s]", aa.newEmailAddress), func() {
						confirmationHash = givenCustomerEmailAddressWasChanged(customerID, aa, 3)

						Convey("When he confirms his changed email address", func() {
							err = ac.confirmCustomerEmailAddress(customerID.String(), confirmationHash.String())
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then his email address should be [%s] and confirmed", aa.newEmailAddress), func() {
								actualCustomerView, err = ac.customerViewByID(customerID.String())
								So(err, ShouldBeNil)
								expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
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

		Convey("\nSCENARIO: A Customer tries to confirm his email address with invalid input", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, confirmationHash = givenCustomerRegistered(aa)

				Convey("When he supplies an empty confirmation hash", func() {
					err = ac.confirmCustomerEmailAddress(customerID.String(), "")

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerEmailAddresses(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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
				customerID, _ = givenCustomerRegistered(aa)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(customerID, aa, 2)

					Convey(fmt.Sprintf("When she changes her email address to [%s]", aa.newEmailAddress), func() {
						err = ac.changeCustomerEmailAddress(customerID.String(), aa.newEmailAddress)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("Then her email address should be [%s] and unconfirmed", aa.newEmailAddress), func() {
							actualCustomerView, err = ac.customerViewByID(customerID.String())
							So(err, ShouldBeNil)
							expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
							expectedCustomerView.EmailAddress = aa.newEmailAddress
							expectedCustomerView.Version = 3
							So(actualCustomerView, ShouldResemble, expectedCustomerView)

							Convey(fmt.Sprintf("And when she tries to change her email address to [%s] again", aa.newEmailAddress), func() {
								err = ac.changeCustomerEmailAddress(customerID.String(), aa.newEmailAddress)
								So(err, ShouldBeNil)

								Convey(fmt.Sprintf("Then her email address should still be [%s]", aa.newEmailAddress), func() {
									actualCustomerView, err = ac.customerViewByID(customerID.String())
									So(err, ShouldBeNil)
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
				customerID, _ = givenCustomerRegistered(aa)

				Convey(fmt.Sprintf("And given she changed her email address to [%s]", aa.newEmailAddress), func() {
					_ = givenCustomerEmailAddressWasChanged(customerID, aa, 2)

					Convey(fmt.Sprintf("And given another Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
						otherCustomerID, _ = givenCustomerRegistered(aa)

						Convey(fmt.Sprintf("When she also tries to change her email address to [%s]", aa.newEmailAddress), func() {
							err = ac.changeCustomerEmailAddress(otherCustomerID.String(), aa.newEmailAddress)

							Convey("Then she should receive an error", func() {
								So(err, ShouldBeError)
								So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to change her email address with invalid input", func() {
			invalidEmailAddress := "fiona@galagher.c"

			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey(fmt.Sprintf("When she supplies an invalid email address [%s]", invalidEmailAddress), func() {
					err = ac.changeCustomerEmailAddress(customerID.String(), invalidEmailAddress)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)

			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerNames(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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
				customerID, _ = givenCustomerRegistered(aa)

				Convey(fmt.Sprintf("When he changes his name to [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
					err = ac.changeCustomerName(customerID.String(), aa.newGivenName, aa.newFamilyName)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then his name should be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
						actualCustomerView, err = ac.customerViewByID(customerID.String())
						So(err, ShouldBeNil)
						expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(customerID, aa)
						expectedCustomerView.GivenName = aa.newGivenName
						expectedCustomerView.FamilyName = aa.newFamilyName
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey(fmt.Sprintf("And when he tries to change his name to [%s %s] again", aa.newGivenName, aa.newFamilyName), func() {
							err = ac.changeCustomerName(customerID.String(), aa.newGivenName, aa.newFamilyName)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then his name should still be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
								actualCustomerView, err = ac.customerViewByID(customerID.String())
								So(err, ShouldBeNil)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to change his name with invalid input", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, _ = givenCustomerRegistered(aa)

				Convey("When he supplies an empty given name", func() {
					err = ac.changeCustomerName(customerID.String(), "", aa.familyName)

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When he supplies an empty family name", func() {
					err = ac.changeCustomerName(customerID.String(), aa.givenName, "")

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForDeletingCustomers(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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
				customerID, confirmationHash = givenCustomerRegistered(aa)

				Convey("When she deletes her account", func() {
					err = ac.deleteCustomer(customerID.String())
					So(err, ShouldBeNil)

					Convey("And when she tries to retrieve her account data", func() {
						actualCustomerView, err = ac.customerViewByID(customerID.String())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})
					})

					Convey("And when she tries to delete her account again", func() {
						err = ac.deleteCustomer(customerID.String())
						So(err, ShouldBeNil)

						Convey("Then her account should still be deleted", func() {
							actualCustomerView, err = ac.customerViewByID(customerID.String())
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})
					})

					Convey("And when she tries to confirm her email address", func() {
						err = ac.confirmCustomerEmailAddress(customerID.String(), confirmationHash.String())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her email address", func() {
						err = ac.changeCustomerEmailAddress(customerID.String(), aa.newEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her name", func() {
						err = ac.changeCustomerName(customerID.String(), aa.newGivenName, aa.newFamilyName)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_WhenCustomerWasNeverRegistered(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

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

		Convey("\nSCENARIO: A hacker tries to play around with a non existing Customer account by guessing IDs", func() {
			Convey("When he tries to retrieve data for a non existing account", func() {
				actualCustomerView, err = ac.customerViewByID(customerID.String())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(actualCustomerView, ShouldBeZeroValue)
				})
			})

			Convey("And when he tries to confirm an email address", func() {
				err = ac.confirmCustomerEmailAddress(customerID.String(), confirmationHash.String())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to change an email address", func() {
				err = ac.changeCustomerEmailAddress(customerID.String(), aa.newEmailAddress)

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to change a name", func() {
				err = ac.changeCustomerName(customerID.String(), aa.newGivenName, aa.newFamilyName)

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when he tries to delete an account", func() {
				err = ac.deleteCustomer(customerID.String())

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}

func TestCustomerAcceptanceScenarios_InvalidClientInput(t *testing.T) {
	ac := bootstrapAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var customerID values.CustomerID
		var confirmationHash values.ConfirmationHash

		aa := acceptanceTestArtifacts{
			emailAddress: "fiona@gallagher.net",
			givenName:    "Fiona",
			familyName:   "Gallagher",
		}

		Convey("\nSCENARIO: A client (web, app, ...) of the application sends an empty id", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				customerID, confirmationHash = givenCustomerRegistered(aa)

				Convey("When she tries to confirm her email address with an empty id", func() {
					err = ac.confirmCustomerEmailAddress("", confirmationHash.String())

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to change her email address with an empty id", func() {
					err = ac.changeCustomerEmailAddress("", aa.emailAddress)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to change her name with an empty id", func() {
					err = ac.changeCustomerName("", aa.givenName, aa.familyName)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to delete her account with an empty id", func() {
					err = ac.deleteCustomer("")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to retrieve her account with an empty id", func() {
					_, err = ac.customerViewByID("")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = acceptanceTestCustomerEventStore.PurgeCustomerEventStream(customerID)
			So(err, ShouldBeNil)
		})
	})
}

func givenCustomerRegistered(
	aa acceptanceTestArtifacts,
) (values.CustomerID, values.ConfirmationHash) {

	customerID := values.GenerateCustomerID()
	emailAddress := values.RebuildEmailAddress(aa.emailAddress)
	confirmationHash := values.GenerateConfirmationHash(emailAddress.String())
	personName := values.RebuildPersonName(aa.givenName, aa.familyName)

	event := events.BuildCustomerRegistered(
		customerID,
		emailAddress,
		confirmationHash,
		personName,
		1,
	)

	err := acceptanceTestCustomerEventStore.RegisterCustomer(es.RecordedEvents{event}, customerID)
	So(err, ShouldBeNil)

	return customerID, confirmationHash
}

func givenCustomerEmailAddressWasConfirmed(
	customerID values.CustomerID,
	aa acceptanceTestArtifacts,
	streamVersion uint,
) {

	emailAddress := values.RebuildEmailAddress(aa.emailAddress)

	event := events.BuildCustomerEmailAddressConfirmed(
		customerID,
		emailAddress,
		streamVersion,
	)

	err := acceptanceTestCustomerEventStore.AppendToCustomerEventStream(es.RecordedEvents{event}, customerID)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	customerID values.CustomerID,
	aa acceptanceTestArtifacts,
	streamVersion uint,
) values.ConfirmationHash {

	emailAddress := values.RebuildEmailAddress(aa.newEmailAddress)
	previousEmailAddress := values.RebuildEmailAddress(aa.emailAddress)
	confirmationHash := values.GenerateConfirmationHash(emailAddress.String())

	event := events.BuildCustomerEmailAddressChanged(
		customerID,
		emailAddress,
		confirmationHash,
		previousEmailAddress,
		streamVersion,
	)

	err := acceptanceTestCustomerEventStore.AppendToCustomerEventStream(es.RecordedEvents{event}, customerID)
	So(err, ShouldBeNil)

	return confirmationHash
}

func bootstrapAcceptanceTestCollaborators() acceptanceTestCollaborators {
	logger := cmd.NewNilLogger()
	config := cmd.MustBuildConfigFromEnv(logger)
	diContainer, err := cmd.Bootstrap(config, logger)
	if err != nil {
		panic(err)
	}

	acceptanceTestCustomerEventStore = diContainer.GetCustomerEventStore()

	return acceptanceTestCollaborators{
		registerCustomer:            diContainer.GetCustomerCommandHandler().RegisterCustomer,
		confirmCustomerEmailAddress: diContainer.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
		changeCustomerEmailAddress:  diContainer.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
		changeCustomerName:          diContainer.GetCustomerCommandHandler().ChangeCustomerName,
		deleteCustomer:              diContainer.GetCustomerCommandHandler().DeleteCustomer,
		customerViewByID:            diContainer.GetCustomerQueryHandler().CustomerViewByID,
	}
}

func buildDefaultCustomerViewForAcceptanceTest(
	customerID values.CustomerID,
	aa acceptanceTestArtifacts,
) customer.View {

	return customer.View{
		ID:                      customerID.String(),
		EmailAddress:            aa.emailAddress,
		IsEmailAddressConfirmed: false,
		GivenName:               aa.givenName,
		FamilyName:              aa.familyName,
		Version:                 1,
	}
}
