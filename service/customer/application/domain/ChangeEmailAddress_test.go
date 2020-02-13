package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		id := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")
		changedEmailAddress := values.RebuildEmailAddress("latoya@ball.net")

		registered := events.ItWasRegistered(
			id,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			registered.CustomerID().ID(),
			changedEmailAddress.EmailAddress(),
		)
		So(err, ShouldBeNil)

		emailAddressConfirmed := events.EmailAddressWasConfirmed(
			registered.CustomerID(),
			registered.EmailAddress(),
			2,
		)

		emailAddressChanged := events.EmailAddressWasChanged(
			registered.CustomerID(),
			registered.EmailAddress(),
			registered.ConfirmationHash(),
			2,
		)

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			changeEmailAddress.CustomerID().ID(),
			changeEmailAddress.EmailAddress().EmailAddress(),
			changeEmailAddress.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's emailAddress", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("When ChangeEmailAddress", func() {
					recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

					Convey("Then EmailAddressChanged", func() {
						ThenEmailAddressChanged(recordedEvents, changeEmailAddress)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's emailAddress twice to an equal value", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("and EmailAddressChanged", func() {
					eventStream = append(eventStream, emailAddressChanged)

					Convey("When ChangeEmailAddress", func() {
						recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Confirm a Customer's changed emailAddress", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("and EmailAddressConfirmed", func() {
					eventStream = append(eventStream, emailAddressConfirmed)

					Convey("and EmailAddressChanged", func() {
						eventStream = append(eventStream, emailAddressChanged)

						Convey("When ConfirmEmailAddress", func() {
							recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

							Convey("Then EmailAddressConfirmed", func() {
								ThenChangedEmailAddressConfirmed(recordedEvents, confirmEmailAddress)
							})
						})
					})
				})
			})
		})
	})

	//Convey("Given a Customer with a confirmed emailAddress", t, func() {
	//	id := values.GenerateCustomerID()
	//	emailAddress := values.RebuildEmailAddress("kevin@ball.com")
	//	confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
	//	personName := values.RebuildPersonName("Kevin", "Ball")
	//
	//	currentStreamVersion := uint(1)
	//	registered := events.ItWasRegistered(id, emailAddress, confirmationHash, personName, currentStreamVersion)
	//
	//	currentStreamVersion++
	//	emailAddressConfirmed := events.EmailAddressWasConfirmed(
	//		registered.CustomerID(),
	//		registered.EmailAddress(),
	//		currentStreamVersion,
	//	)
	//
	//	eventStream := es.DomainEvents{registered, emailAddressConfirmed}
	//
	//	Convey("When the emailAddress is changed", func() {
	//		changeEmailAddress, err := commands.NewChangeEmailAddress(
	//			registered.CustomerID().ID(),
	//			"john+changed@doe.com",
	//		)
	//		So(err, ShouldBeNil)
	//
	//		recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)
	//
	//		Convey("Then it should record EmailAddressChanged", func() {
	//			So(recordedEvents, ShouldHaveLength, 1)
	//			emailAddressChanged, ok := recordedEvents[0].(events.EmailAddressChanged)
	//			So(ok, ShouldBeTrue)
	//			So(emailAddressChanged, ShouldNotBeNil)
	//			So(emailAddressChanged.CustomerID().Equals(changeEmailAddress.CustomerID()), ShouldBeTrue)
	//			So(emailAddressChanged.EmailAddress().Equals(changeEmailAddress.EmailAddress()), ShouldBeTrue)
	//			So(emailAddressChanged.ConfirmationHash().Equals(changeEmailAddress.ConfirmationHash()), ShouldBeTrue)
	//			So(emailAddressChanged.StreamVersion(), ShouldEqual, currentStreamVersion+1)
	//
	//			eventStream = append(eventStream, emailAddressChanged)
	//			currentStreamVersion++
	//
	//			Convey("And when it is changed to the same value again", func() {
	//				recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)
	//
	//				Convey("Then it should be ignored", func() {
	//					So(recordedEvents, ShouldBeEmpty)
	//				})
	//			})
	//
	//			Convey("And when the changed emailAddress is confirmed with the right confirmationHash", func() {
	//				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
	//					emailAddressChanged.CustomerID().ID(),
	//					emailAddressChanged.EmailAddress().EmailAddress(),
	//					emailAddressChanged.ConfirmationHash().Hash(),
	//				)
	//				So(err, ShouldBeNil)
	//
	//				recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)
	//
	//				Convey("Then it should record EmailAddressConfirmed", func() {
	//					So(recordedEvents, ShouldHaveLength, 1)
	//					emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
	//					So(ok, ShouldBeTrue)
	//					So(emailAddressConfirmed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
	//					So(emailAddressConfirmed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
	//					So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)
	//				})
	//			})
	//		})
	//	})
	//})
}

func ThenEmailAddressChanged(
	recordedEvents es.DomainEvents,
	changeEmailAddress commands.ChangeEmailAddress,
) {

	So(recordedEvents, ShouldHaveLength, 1)
	emailAddressChanged, ok := recordedEvents[0].(events.EmailAddressChanged)
	So(ok, ShouldBeTrue)
	So(emailAddressChanged, ShouldNotBeNil)
	So(emailAddressChanged.CustomerID().Equals(changeEmailAddress.CustomerID()), ShouldBeTrue)
	So(emailAddressChanged.EmailAddress().Equals(changeEmailAddress.EmailAddress()), ShouldBeTrue)
	So(emailAddressChanged.ConfirmationHash().Equals(changeEmailAddress.ConfirmationHash()), ShouldBeTrue)
	So(emailAddressChanged.StreamVersion(), ShouldEqual, 2)
}

func ThenChangedEmailAddressConfirmed(
	recordedEvents es.DomainEvents,
	confirmEmailAddress commands.ConfirmEmailAddress,
) {

	So(recordedEvents, ShouldHaveLength, 1)
	emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
	So(ok, ShouldBeTrue)
	So(emailAddressConfirmed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
	So(emailAddressConfirmed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
	So(emailAddressConfirmed.StreamVersion(), ShouldEqual, 4)
}
