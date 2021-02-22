package customeraccounts_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var atStartCustomerEventStream application.ForStartingCustomerEventStreams
var atAppendToCustomerEventStream application.ForAppendingToCustomerEventStreams
var atPurgeCustomerEventStream application.ForPurgingCustomerEventStreams

type acceptanceTestCollaborators struct {
	registerCustomer            hexagon.ForRegisteringCustomers
	confirmCustomerEmailAddress hexagon.ForConfirmingCustomerEmailAddresses
	changeCustomerEmailAddress  hexagon.ForChangingCustomerEmailAddresses
	changeCustomerName          hexagon.ForChangingCustomerNames
	deleteCustomer              hexagon.ForDeletingCustomers
	customerViewByID            hexagon.ForRetrievingCustomerViews
}

type acceptanceTestValues struct {
	customerID          value.CustomerID
	otherCustomerID     value.CustomerID
	emailAddress        value.UnconfirmedEmailAddress
	changedEmailAddress value.UnconfirmedEmailAddress
	name                value.PersonName
	changedName         value.PersonName
	ea                  string // emailAddress
	cea                 string // changeEmailAddress
	ch                  string // confirmationHash
	cch                 string // changeConfirmationHash
	gn                  string // givenName
	fn                  string // familyName
	cgn                 string // changeGivenName
	cfn                 string // changedFamilyName
}

func TestCustomerAcceptanceScenarios_ForRegisteringCustomers(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A prospective Customer registers her account", func() {
			Convey(fmt.Sprintf("When a Customer registers as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				err = ac.registerCustomer(v.customerID, v.ea, v.gn, v.fn)
				So(err, ShouldBeNil)

				expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
				details := fmt.Sprintf("\n\tGivenName: %s", expectedCustomerView.GivenName)
				details += fmt.Sprintf("\n\tFamilyName: %s", expectedCustomerView.FamilyName)
				details += fmt.Sprintf("\n\tEmailAddress: %s", expectedCustomerView.EmailAddress)
				details += fmt.Sprintf("\n\tIsEmailAddressConfirmed: %t", expectedCustomerView.IsEmailAddressConfirmed)

				Convey(fmt.Sprintf("Then her account should show the data she supplied: %s", details), func() {
					actualCustomerView, err = ac.customerViewByID(v.customerID.String())
					So(err, ShouldBeNil)
					So(actualCustomerView, ShouldResemble, expectedCustomerView)
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can't register because her email address is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", v.ea), func() {
					err = ac.registerCustomer(v.customerID, v.ea, v.gn, v.fn)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can register with an email address that is not used any more", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("And given the first Customer deleted her account", func() {
					err = ac.deleteCustomer(v.customerID.String())
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", v.ea), func() {
						err = ac.registerCustomer(v.otherCustomerID, v.ea, v.gn, v.fn)

						Convey("Then she should be able to register", func() {
							So(err, ShouldBeNil)
						})
					})
				})

				Convey(fmt.Sprintf("Or given the first Customer changed her email address to [%s]", v.cea), func() {
					err = ac.changeCustomerEmailAddress(v.customerID.String(), v.cea)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", v.ea), func() {
						err = ac.registerCustomer(v.otherCustomerID, v.ea, v.gn, v.fn)

						Convey("Then she should be able to register", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer tries to register with invalid input", func() {
			invalidEmailAddress := "fiona@galagher.c"

			Convey(fmt.Sprintf("When she supplies an invalid email address [%s]", invalidEmailAddress), func() {
				err = ac.registerCustomer(v.customerID, invalidEmailAddress, v.gn, v.fn)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she supplies an empty givenName", func() {
				err = ac.registerCustomer(v.customerID, v.ea, "", v.fn)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When she supplies an empty familyName", func() {
				err = ac.registerCustomer(v.customerID, v.ea, v.gn, "")

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)

			err = atPurgeCustomerEventStream(v.otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForConfirmingCustomerEmailAddresses(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer confirms her email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she confirms her email address", func() {
					err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.ch)
					So(err, ShouldBeNil)

					Convey("Then her email address should be confirmed", func() {
						actualCustomerView, err = ac.customerViewByID(v.customerID.String())
						So(err, ShouldBeNil)
						expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
						expectedCustomerView.IsEmailAddressConfirmed = true
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey("And when she confirms her email address again", func() {
							err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.ch)
							So(err, ShouldBeNil)

							Convey("Then her email address should still be confirmed", func() {
								actualCustomerView, err = ac.customerViewByID(v.customerID.String())
								So(err, ShouldBeNil)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer can't confirm her email address, because the confirmation hash is not matching", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she tries to confirm her email address with a wrong confirmation hash", func() {
					err = ac.confirmCustomerEmailAddress(v.customerID.String(), "invalid_confirmation_hash")

					Convey("Then she should receive an error", func() {
						So(errors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)

						Convey("And her email address should still be unconfirmed", func() {
							actualCustomerView, err = ac.customerViewByID(v.customerID.String())
							So(err, ShouldBeNil)
							expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her already confirmed email address again", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(v.customerID, v.emailAddress, 2)

					Convey("When she tries to confirm her email address again with a wrong confirmation hash", func() {
						err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.ch)
						So(err, ShouldBeNil)

						Convey("Then her email address should still be confirmed", func() {
							actualCustomerView, err = ac.customerViewByID(v.customerID.String())
							So(err, ShouldBeNil)

							expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
							expectedCustomerView.IsEmailAddressConfirmed = true
							expectedCustomerView.Version = 2
							So(actualCustomerView, ShouldResemble, expectedCustomerView)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer confirms her changed email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(v.customerID, v.emailAddress, 2)

					Convey(fmt.Sprintf("And given she changed her email address to [%s]", v.cea), func() {
						givenCustomerEmailAddressWasChanged(v.customerID, v.changedEmailAddress, 3)

						Convey("When she confirms her changed email address", func() {
							err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.cch)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her email address should be [%s] and confirmed", v.cea), func() {
								actualCustomerView, err = ac.customerViewByID(v.customerID.String())
								So(err, ShouldBeNil)
								expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
								expectedCustomerView.EmailAddress = v.cea
								expectedCustomerView.IsEmailAddressConfirmed = true
								expectedCustomerView.Version = 4
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to confirm her email address with invalid input", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she supplies an empty confirmation hash", func() {
					err = ac.confirmCustomerEmailAddress(v.customerID.String(), "")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerEmailAddresses(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer changes her (confirmed) email address", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("And given she confirmed her email address", func() {
					givenCustomerEmailAddressWasConfirmed(v.customerID, v.emailAddress, 2)

					Convey(fmt.Sprintf("When she changes her email address to [%s]", v.cea), func() {
						err = ac.changeCustomerEmailAddress(v.customerID.String(), v.cea)
						So(err, ShouldBeNil)

						Convey(fmt.Sprintf("Then her email address should be [%s] and unconfirmed", v.cea), func() {
							actualCustomerView, err = ac.customerViewByID(v.customerID.String())
							So(err, ShouldBeNil)
							expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
							expectedCustomerView.EmailAddress = v.cea
							expectedCustomerView.Version = 3
							So(actualCustomerView, ShouldResemble, expectedCustomerView)

							Convey(fmt.Sprintf("And when she tries to change her email address to [%s] again", v.cea), func() {
								err = ac.changeCustomerEmailAddress(v.customerID.String(), v.cea)
								So(err, ShouldBeNil)

								Convey(fmt.Sprintf("Then her email address should still be [%s]", v.cea), func() {
									actualCustomerView, err = ac.customerViewByID(v.customerID.String())
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
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey(fmt.Sprintf("And given she changed her email address to [%s]", v.cea), func() {
					givenCustomerEmailAddressWasChanged(v.customerID, v.changedEmailAddress, 2)

					Convey(fmt.Sprintf("And given another Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
						givenCustomerRegistered(v.otherCustomerID, v.emailAddress, v.name)

						Convey(fmt.Sprintf("When she also tries to change her email address to [%s]", v.cea), func() {
							err = ac.changeCustomerEmailAddress(v.otherCustomerID.String(), v.cea)

							Convey("Then she should receive an error", func() {
								So(err, ShouldBeError)
								So(errors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to change her email address with invalid input", func() {
			invalidEmailAddress := "fiona@galagher.c"

			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey(fmt.Sprintf("When she supplies an invalid email address [%s]", invalidEmailAddress), func() {
					err = ac.changeCustomerEmailAddress(v.customerID.String(), invalidEmailAddress)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)

			err = atPurgeCustomerEventStream(v.otherCustomerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForChangingCustomerNames(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var expectedCustomerView customer.View
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer changes her name", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey(fmt.Sprintf("When she changes her name to [%s %s]", v.cgn, v.cfn), func() {
					err = ac.changeCustomerName(v.customerID.String(), v.cgn, v.cfn)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then her name should be [%s %s]", v.cgn, v.cfn), func() {
						actualCustomerView, err = ac.customerViewByID(v.customerID.String())
						So(err, ShouldBeNil)
						expectedCustomerView = buildDefaultCustomerViewForAcceptanceTest(v.customerID, v.emailAddress, v.name)
						expectedCustomerView.GivenName = v.cgn
						expectedCustomerView.FamilyName = v.cfn
						expectedCustomerView.Version = 2
						So(actualCustomerView, ShouldResemble, expectedCustomerView)

						Convey(fmt.Sprintf("And when she tries to change her name to [%s %s] again", v.cgn, v.cfn), func() {
							err = ac.changeCustomerName(v.customerID.String(), v.cgn, v.cfn)
							So(err, ShouldBeNil)

							Convey(fmt.Sprintf("Then her name should still be [%s %s]", v.cgn, v.cfn), func() {
								actualCustomerView, err = ac.customerViewByID(v.customerID.String())
								So(err, ShouldBeNil)
								So(actualCustomerView, ShouldResemble, expectedCustomerView)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to change her name with invalid input", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she supplies an empty given name", func() {
					err = ac.changeCustomerName(v.customerID.String(), "", v.cfn)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she supplies an empty family name", func() {
					err = ac.changeCustomerName(v.customerID.String(), v.cgn, "")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_ForAddingBillingProfiles(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer adds a billing profile to her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				Convey("When she adds a billing profile ", func() {
					Convey("Then her account should contain one billing profile", func() {

					})
				})
			})
		})
	})
}

func TestCustomerAcceptanceScenarios_ForRemovingBillingProfiles(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer removes a billing profile from her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				Convey("And given she added a billing profile", func() {
					Convey("And given she added another billing profile", func() {
						Convey("When she removes the first billing profile from her account", func() {
							Convey("Then her account should contain only the second billing profile", func() {

							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: A Customer tries to remove the only billing profile from her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				Convey("And given she added a billing profile", func() {
					Convey("When she tries to remove the billing profile from her account", func() {
						Convey("Then she should receive an error", func() {

						})
					})
				})
			})
		})
	})
}

func TestCustomerAcceptanceScenarios_ForDeletingCustomers(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A Customer deletes her account", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she deletes her account", func() {
					err = ac.deleteCustomer(v.customerID.String())
					So(err, ShouldBeNil)

					Convey("And when she tries to retrieve her account data", func() {
						actualCustomerView, err = ac.customerViewByID(v.customerID.String())

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})
					})

					Convey("And when she tries to delete her account again", func() {
						err = ac.deleteCustomer(v.customerID.String())
						So(err, ShouldBeNil)

						Convey("Then her account should still be deleted", func() {
							actualCustomerView, err = ac.customerViewByID(v.customerID.String())
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
							So(actualCustomerView, ShouldBeZeroValue)
						})
					})

					Convey("And when she tries to confirm her email address", func() {
						err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.ch)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her email address", func() {
						err = ac.changeCustomerEmailAddress(v.customerID.String(), v.cea)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
						})
					})

					Convey("And when she tries to change her name", func() {
						err = ac.changeCustomerName(v.customerID.String(), v.cgn, v.cfn)

						Convey("Then she should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)
		})
	})
}

func TestCustomerAcceptanceScenarios_WhenCustomerWasNeverRegistered(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error
		var actualCustomerView customer.View

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A hacker tries to play around with a non existing Customer account by guessing IDs", func() {
			Convey("When she tries to retrieve data for a non existing account", func() {
				actualCustomerView, err = ac.customerViewByID(v.customerID.String())

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
					So(actualCustomerView, ShouldBeZeroValue)
				})
			})

			Convey("And when she tries to confirm an email address", func() {
				err = ac.confirmCustomerEmailAddress(v.customerID.String(), v.ch)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when she tries to change an email address", func() {
				err = ac.changeCustomerEmailAddress(v.customerID.String(), v.ea)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when she tries to change a name", func() {
				err = ac.changeCustomerName(v.customerID.String(), v.gn, v.fn)

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				})
			})

			Convey("And when she tries to delete an account", func() {
				err = ac.deleteCustomer(v.customerID.String())

				Convey("Then she should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}

func TestCustomerAcceptanceScenarios_InvalidClientInput(t *testing.T) {
	ac := initAcceptanceTestCollaborators()

	Convey("Prepare test artifacts", t, func() {
		var err error

		v := initAcceptanceTestValues()

		Convey("\nSCENARIO: A client (web, app, ...) of the application sends an empty id", func() {
			Convey(fmt.Sprintf("Given a Customer registered as [%s %s] with [%s]", v.gn, v.fn, v.ea), func() {
				givenCustomerRegistered(v.customerID, v.emailAddress, v.name)

				Convey("When she tries to confirm her email address with an empty id", func() {
					err = ac.confirmCustomerEmailAddress("", v.ch)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to change her email address with an empty id", func() {
					err = ac.changeCustomerEmailAddress("", v.ea)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to change her name with an empty id", func() {
					err = ac.changeCustomerName("", v.gn, v.fn)

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to delete her account with an empty id", func() {
					err = ac.deleteCustomer("")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})

				Convey("When she tries to retrieve her account with an empty id", func() {
					_, err = ac.customerViewByID("")

					Convey("Then she should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() {
			err = atPurgeCustomerEventStream(v.customerID)
			So(err, ShouldBeNil)
		})
	})
}

func givenCustomerRegistered(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	name value.PersonName,
) {

	registered := domain.BuildCustomerRegistered(
		customerID,
		emailAddress,
		name,
		es.GenerateMessageID(),
		1,
	)

	err := atStartCustomerEventStream(registered)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasConfirmed(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	streamVersion uint,
) {

	confirmedEmailAddress, err := value.ConfirmEmailAddressWithHash(emailAddress, emailAddress.ConfirmationHash())
	So(err, ShouldBeNil)

	event := domain.BuildCustomerEmailAddressConfirmed(
		customerID,
		confirmedEmailAddress,
		es.GenerateMessageID(),
		streamVersion,
	)

	err = atAppendToCustomerEventStream(es.RecordedEvents{event}, customerID)
	So(err, ShouldBeNil)
}

func givenCustomerEmailAddressWasChanged(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	streamVersion uint,
) {

	event := domain.BuildCustomerEmailAddressChanged(
		customerID,
		emailAddress,
		es.GenerateMessageID(),
		streamVersion,
	)

	err := atAppendToCustomerEventStream(es.RecordedEvents{event}, customerID)
	So(err, ShouldBeNil)
}

func buildDefaultCustomerViewForAcceptanceTest(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	name value.PersonName,
) customer.View {

	return customer.View{
		ID:                      customerID.String(),
		EmailAddress:            emailAddress.String(),
		IsEmailAddressConfirmed: false,
		GivenName:               name.GivenName(),
		FamilyName:              name.FamilyName(),
		Version:                 1,
	}
}

func initAcceptanceTestCollaborators() acceptanceTestCollaborators {
	logger := shared.NewNilLogger()
	config := grpc.MustBuildConfigFromEnv(logger)
	postgresDBConn := grpc.MustInitPostgresDB(config, logger)
	diContainer := grpc.MustBuildDIContainer(config, logger, grpc.UsePostgresDBConn(postgresDBConn))
	eventStore := diContainer.GetCustomerEventStore()
	atStartCustomerEventStream = eventStore.StartEventStream
	atAppendToCustomerEventStream = eventStore.AppendToEventStream
	atPurgeCustomerEventStream = eventStore.PurgeEventStream

	return acceptanceTestCollaborators{
		registerCustomer:            diContainer.GetCustomerCommandHandler().RegisterCustomer,
		confirmCustomerEmailAddress: diContainer.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
		changeCustomerEmailAddress:  diContainer.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
		changeCustomerName:          diContainer.GetCustomerCommandHandler().ChangeCustomerName,
		deleteCustomer:              diContainer.GetCustomerCommandHandler().DeleteCustomer,
		customerViewByID:            diContainer.GetCustomerQueryHandler().CustomerViewByID,
	}
}

func initAcceptanceTestValues() acceptanceTestValues {
	customerID := value.GenerateCustomerID()
	otherCustomerID := value.GenerateCustomerID()
	emailAddress, err := value.BuildUnconfirmedEmailAddress("fiona@gallagher.net")
	So(err, ShouldBeNil)
	changedEmailAddress, err := value.BuildUnconfirmedEmailAddress("fiona@lishman.net")
	So(err, ShouldBeNil)
	name, err := value.BuildPersonName("Fiona", "Gallagher")
	So(err, ShouldBeNil)
	changedName, err := value.BuildPersonName("Fiona", "Lishman")
	So(err, ShouldBeNil)

	return acceptanceTestValues{
		customerID:          customerID,
		otherCustomerID:     otherCustomerID,
		emailAddress:        emailAddress,
		changedEmailAddress: changedEmailAddress,
		name:                name,
		changedName:         changedName,
		ea:                  emailAddress.String(),
		cea:                 changedEmailAddress.String(),
		ch:                  emailAddress.ConfirmationHash().String(),
		cch:                 changedEmailAddress.ConfirmationHash().String(),
		gn:                  name.GivenName(),
		fn:                  name.FamilyName(),
		cgn:                 changedName.GivenName(),
		cfn:                 changedName.FamilyName(),
	}
}
