package values_test

import (
	"fmt"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

/*** Tests for Getter methods ***/

func TestConfirmableEmailAddressExposesExpectedValues(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.EmailAddressFrom(emailAddressValue)
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		Convey("It should expose the expected EmailAddress", func() {
			So(confirmableEmailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})

		Convey("And it should expose a ConfirmationHash", func() {
			So(confirmableEmailAddress.ConfirmationHash(), ShouldNotBeBlank)
		})

		Convey("And it should not be confirmed", func() {
			So(confirmableEmailAddress.IsConfirmed(), ShouldBeFalse)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestConfirmableEmailAddressEquals(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.EmailAddressFrom(emailAddressValue)
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		Convey("And given an equal EmailAddress", func() {
			equalConfirmableEmailAddress, err := values.EmailAddressFrom(emailAddressValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := confirmableEmailAddress.Equals(equalConfirmableEmailAddress)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given a different EmailAddress", func() {
			differentEmailAddress, err := values.EmailAddressFrom("foo+different@bar.com")
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := confirmableEmailAddress.Equals(differentEmailAddress)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}

func TestConfirmableEmailAddressMarkAsConfirmed(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddress, err := values.EmailAddressFrom("foo@bar.com")
		So(err, ShouldBeNil)
		unconfirmedEmailAddress := emailAddress.ToConfirmable()

		Convey("When it is marked as confirmed", func() {
			confirmedEmailAddress := unconfirmedEmailAddress.MarkAsConfirmed()

			Convey("It should be confirmed", func() {
				So(confirmedEmailAddress.IsConfirmed(), ShouldBeTrue)
				So(unconfirmedEmailAddress.IsConfirmed(), ShouldBeFalse)
			})
		})
	})
}

/*** Tests for Modification methods ***/

func TestConfirmableEmailAddressIsConfirmedBy(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddress, err := values.EmailAddressFrom("foo@bar.com")
		So(err, ShouldBeNil)
		unconfirmedEmailAddress := emailAddress.ToConfirmable()

		Convey("It should confirm with the matching ConfirmationHash", func() {
			confirmationHash, err := values.ConfirmationHashFrom(unconfirmedEmailAddress.ConfirmationHash())
			So(err, ShouldBeNil)
			So(unconfirmedEmailAddress.IsConfirmedBy(confirmationHash), ShouldBeTrue)
		})

		Convey("It should not confirm with a wrong ConfirmationHash", func() {
			confirmationHash, err := values.ConfirmationHashFrom("invalid_confirmation_hash")
			So(err, ShouldBeNil)
			So(unconfirmedEmailAddress.IsConfirmedBy(confirmationHash), ShouldBeFalse)
		})
	})
}

/*** Tests for Marshal/Unmarshal methods ***/

func TestConfirmableEmailAddressMarshalJSON(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddress, err := values.EmailAddressFrom("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		Convey("When it is marshaled to json", func() {
			data, err := confirmableEmailAddress.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)

				expectedJSON := fmt.Sprintf(
					`{"emailAddress":"%s","confirmationHash":"%s"}`,
					confirmableEmailAddress.EmailAddress(),
					confirmableEmailAddress.ConfirmationHash(),
				)

				So(string(data), ShouldEqual, expectedJSON)
			})
		})
	})
}

func TestConfirmableEmailAddressUnmarshalJSON(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress marshaled to json", t, func() {
		emailAddress, err := values.EmailAddressFrom("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		data, err := confirmableEmailAddress.MarshalJSON()

		Convey("And when it is unmarshaled", func() {
			unmarshaled := &values.ConfirmableEmailAddress{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original ConfirmableEmailAddress", func() {
				So(err, ShouldBeNil)
				So(confirmableEmailAddress, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to ConfirmableEmailAddress", func() {
			unmarshaled := &values.ConfirmableEmailAddress{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
