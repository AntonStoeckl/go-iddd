package application_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/query"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerScenarios(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()
	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var actualCustomerView customer.View

		aa := buildArtifactsForAcceptanceTest()

		Convey("\nSCENARIO: A prospective Customer registers her account", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				err = commandHandler.RegisterCustomer(aa.registerCustomer)
				So(err, ShouldBeNil)

				details := fmt.Sprintf("\n\tGivenName: %s", aa.givenName)
				details += fmt.Sprintf("\n\tFamilyName: %s", aa.familyName)
				details += fmt.Sprintf("\n\tEmailAddress: %s (unconfirmed)", aa.emailAddress)

				Convey(fmt.Sprintf("Then her account should show the data she supplied: %s", details), func() {
					actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
					So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can't register because email address is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
					err = commandHandler.RegisterCustomer(aa.registerDuplicateCustomer)

					Convey(fmt.Sprintf("Then she should receive an error"), func() {
						So(err, ShouldBeError)
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can register with an email address that is not used any more", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("And given the first Customer deleted her account", func() {
					err = commandHandler.DeleteCustomer(aa.deleteCustomer)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						err = commandHandler.RegisterCustomer(aa.registerDuplicateCustomer)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})

				Convey("And given the first Customer changed her email address", func() {
					err = commandHandler.ChangeCustomerEmailAddress(aa.changeCustomerEmailAddress)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", aa.emailAddress), func() {
						err = commandHandler.RegisterCustomer(aa.registerDuplicateCustomer)

						Convey(fmt.Sprintf("Then she should be able to register"), func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})

			Reset(func() {
				err = diContainer.GetCustomerEventStore().Delete(aa.deleteDuplicateCustomer.CustomerID())
				So(err, ShouldBeNil)
			})
		})

		Convey("\nSCENARIO: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("When she confirms her email address", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmCustomerEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then her email address should be confirmed", func() {
						actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
						aa.expectedCustomerView.IsEmailAddressConfirmed = true
						aa.expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)

						Convey("And when she confirms her email address again", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmCustomerEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then her email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
								aa.expectedCustomerView.IsEmailAddressConfirmed = true
								aa.expectedCustomerView.Version = 2
								So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't confirm her email address, because the confirmation hash is invalid", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("When she tries to confirm her email address with a wrong confirmation hash", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmCustomerEmailAddressWithInvalidHash)

					Convey("Then she should receive an error", func() {
						So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("And her email address should still be unconfirmed", func() {
							actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
							aa.expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer fails to confirm her already confirmed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(aa.confirmCustomerEmailAddress, commandHandler)

					Convey("When she tries to confirm her email address again with a wrong confirmation hash", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmCustomerEmailAddressWithInvalidHash)

						Convey("Then she should receive an error", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)

							Convey("And her email address should still be confirmed", func() {
								actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
								aa.expectedCustomerView.IsEmailAddressConfirmed = true
								aa.expectedCustomerView.Version = 3
								So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer changes her (confirmed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(aa.confirmCustomerEmailAddress, commandHandler)

					Convey(fmt.Sprintf("When she changes her email address to [%s]", aa.newEmailAddress), func() {
						err = commandHandler.ChangeCustomerEmailAddress(aa.changeCustomerEmailAddress)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("Then her email address should be [%s] and unconfirmed", aa.newEmailAddress), func() {
							actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
							aa.expectedCustomerView.EmailAddress = aa.newEmailAddress
							aa.expectedCustomerView.Version = 3
							So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)

							Convey(fmt.Sprintf("And when she tries to change her email address to [%s] again", aa.newEmailAddress), func() {
								err = commandHandler.ChangeCustomerEmailAddress(aa.changeCustomerEmailAddress)
								So(err, ShouldBeNil)

								Convey(fmt.Sprintf("Then her email address should still be [%s]", aa.newEmailAddress), func() {
									actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
									aa.expectedCustomerView.EmailAddress = aa.newEmailAddress
									aa.expectedCustomerView.Version = 3
									So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
								})
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her (changed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(aa.confirmCustomerEmailAddress, commandHandler)

					Convey(fmt.Sprintf("And given she changed her email address to [%s]", aa.newEmailAddress), func() {
						givenCustomerEmailAddressWasChanged(aa.changeCustomerEmailAddress, commandHandler)

						Convey("When she confirms her changed email address", func() {
							err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmChangedCustomerEmailAddress)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her email address should be [%s] and confirmed", aa.newEmailAddress), func() {
								actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
								aa.expectedCustomerView.EmailAddress = aa.newEmailAddress
								aa.expectedCustomerView.IsEmailAddressConfirmed = true
								aa.expectedCustomerView.Version = 4
								So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't change her email address because it is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey(fmt.Sprintf("And given she changed here email address to [%s]", aa.newEmailAddress), func() {
					givenCustomerEmailAddressWasChanged(aa.changeCustomerEmailAddress, commandHandler)

					Convey(fmt.Sprintf("And given another Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
						givenCustomerRegistered(aa.registerDuplicateCustomer, commandHandler)

						Convey(fmt.Sprintf("When she also tries to change her email address to [%s]", aa.newEmailAddress), func() {
							var changeEmailAddressToAnAlreadyUsedOne commands.ChangeCustomerEmailAddress

							changeEmailAddressToAnAlreadyUsedOne, err = commands.BuildChangeCustomerEmailAddress(
								aa.registerDuplicateCustomer.CustomerID().ID(),
								aa.newEmailAddress,
							)
							So(err, ShouldBeNil)

							err = commandHandler.ChangeCustomerEmailAddress(changeEmailAddressToAnAlreadyUsedOne)

							Convey("Then she should receive an error", func() {
								So(err, ShouldBeError)
								So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
							})
						})
					})
				})
			})

			Reset(func() {
				err = diContainer.GetCustomerEventStore().Delete(aa.deleteDuplicateCustomer.CustomerID())
				So(err, ShouldBeNil)
			})
		})

		Convey("\nSCENARIO: A Customer changes her name", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey(fmt.Sprintf("When she changes her name to [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
					err = commandHandler.ChangeCustomerName(aa.changeCustomerName)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then her name should be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
						actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
						aa.expectedCustomerView.GivenName = aa.newGivenName
						aa.expectedCustomerView.FamilyName = aa.newFamilyName
						aa.expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)

						Convey(fmt.Sprintf("And when she tries to change her name to [%s %s] again", aa.newGivenName, aa.newFamilyName), func() {
							err = commandHandler.ChangeCustomerName(aa.changeCustomerName)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her name should still be [%s %s]", aa.newGivenName, aa.newFamilyName), func() {
								actualCustomerView = retrieveAccountData(queryHandler, aa.customerID)
								aa.expectedCustomerView.GivenName = aa.newGivenName
								aa.expectedCustomerView.FamilyName = aa.newFamilyName
								aa.expectedCustomerView.Version = 2
								So(actualCustomerView, ShouldResemble, aa.expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer deletes her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", aa.givenName, aa.familyName, aa.emailAddress), func() {
				givenCustomerRegistered(aa.registerCustomer, commandHandler)

				Convey("When she deletes her account", func() {
					err = commandHandler.DeleteCustomer(aa.deleteCustomer)
					So(err, ShouldBeNil)

					Convey("And when she tries to retrieve her account data", func() {
						actualCustomerView, err = queryHandler.CustomerViewByID(aa.customerID.ID())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})

						Convey("And when she tries to delete her account again", func() {
							err = commandHandler.DeleteCustomer(aa.deleteCustomer)
							So(err, ShouldBeNil)

							Convey("Then her account should still be deleted", func() {
								actualCustomerView, err = queryHandler.CustomerViewByID(aa.customerID.ID())
								So(err, ShouldBeError)
								So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
								So(actualCustomerView, ShouldBeZeroValue)
							})
						})
					})

					Convey("And when she tries to confirm her email address", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(aa.confirmCustomerEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her email address", func() {
						err = commandHandler.ChangeCustomerEmailAddress(aa.changeCustomerEmailAddress)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her name", func() {
						err = commandHandler.ChangeCustomerName(aa.changeCustomerName)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err = diContainer.GetCustomerEventStore().Delete(aa.customerID)
			So(err, ShouldBeNil)
		})
	})
}

func givenCustomerRegistered(
	registerCustomer commands.RegisterCustomer,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.RegisterCustomer(registerCustomer)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasConfirmed(
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	changeCustomerEmailAddress commands.ChangeCustomerEmailAddress,
	commandHandler *command.CustomerCommandHandler,
) {

	err := commandHandler.ChangeCustomerEmailAddress(changeCustomerEmailAddress)
	So(err, ShouldBeNil)
}

func retrieveAccountData(
	queryHandler *query.CustomerQueryHandler,
	id values.CustomerID,
) customer.View {

	customerView, err := queryHandler.CustomerViewByID(id.ID())
	So(err, ShouldBeNil)

	return customerView
}

type acceptanceTestArtifacts struct {
	emailAddress                               string
	givenName                                  string
	familyName                                 string
	newEmailAddress                            string
	newGivenName                               string
	newFamilyName                              string
	customerID                                 values.CustomerID
	registerCustomer                           commands.RegisterCustomer
	registerDuplicateCustomer                  commands.RegisterCustomer
	confirmCustomerEmailAddress                commands.ConfirmCustomerEmailAddress
	confirmCustomerEmailAddressWithInvalidHash commands.ConfirmCustomerEmailAddress
	confirmChangedCustomerEmailAddress         commands.ConfirmCustomerEmailAddress
	changeCustomerEmailAddress                 commands.ChangeCustomerEmailAddress
	changeCustomerName                         commands.ChangeCustomerName
	deleteCustomer                             commands.DeleteCustomer
	deleteDuplicateCustomer                    commands.DeleteCustomer
	expectedCustomerView                       customer.View
}

func buildArtifactsForAcceptanceTest() acceptanceTestArtifacts {
	var err error
	var aa acceptanceTestArtifacts

	aa.emailAddress = "fiona@gallagher.net"
	aa.givenName = "Fiona"
	aa.familyName = "Galagher"
	aa.newEmailAddress = "fiona@pratt.net"
	aa.newGivenName = "Fiona"
	aa.newFamilyName = "Pratt"

	aa.registerCustomer, err = commands.BuildRegisterCustomer(
		aa.emailAddress,
		aa.givenName,
		aa.familyName,
	)
	So(err, ShouldBeNil)

	aa.customerID = aa.registerCustomer.CustomerID()

	aa.registerDuplicateCustomer, err = commands.BuildRegisterCustomer(
		aa.emailAddress,
		aa.givenName,
		aa.familyName,
	)
	So(err, ShouldBeNil)

	aa.confirmCustomerEmailAddress, err = commands.BuildConfirmCustomerEmailAddress(
		aa.registerCustomer.CustomerID().ID(),
		aa.registerCustomer.ConfirmationHash().Hash(),
	)
	So(err, ShouldBeNil)

	aa.confirmCustomerEmailAddressWithInvalidHash, err = commands.BuildConfirmCustomerEmailAddress(
		aa.registerCustomer.CustomerID().ID(),
		values.GenerateConfirmationHash(aa.emailAddress).Hash(),
	)
	So(err, ShouldBeNil)

	aa.changeCustomerEmailAddress, err = commands.BuildChangeCustomerEmailAddress(
		aa.registerCustomer.CustomerID().ID(),
		aa.newEmailAddress,
	)
	So(err, ShouldBeNil)

	aa.confirmChangedCustomerEmailAddress, err = commands.BuildConfirmCustomerEmailAddress(
		aa.registerCustomer.CustomerID().ID(),
		aa.changeCustomerEmailAddress.ConfirmationHash().Hash(),
	)
	So(err, ShouldBeNil)

	aa.changeCustomerName, err = commands.BuildChangeCustomerName(
		aa.registerCustomer.CustomerID().ID(),
		aa.newGivenName,
		aa.newFamilyName,
	)
	So(err, ShouldBeNil)

	aa.deleteCustomer, err = commands.BuildDeleteCustomer(aa.customerID.ID())
	So(err, ShouldBeNil)

	aa.deleteDuplicateCustomer, err = commands.BuildDeleteCustomer(aa.registerDuplicateCustomer.CustomerID().ID())
	So(err, ShouldBeNil)

	aa.expectedCustomerView = customer.View{
		ID:                      aa.customerID.ID(),
		EmailAddress:            aa.emailAddress,
		IsEmailAddressConfirmed: false,
		GivenName:               aa.givenName,
		FamilyName:              aa.familyName,
		Version:                 1,
	}

	return aa
}
