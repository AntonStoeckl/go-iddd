package values_test

import (
	"fmt"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

/*** Tests for Getter methods ***/

func TestConfirmableEmailAddressExposesExpectedValues(t *testing.T) {
	Convey("Given a new ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
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
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()

		Convey("And given an equal EmailAddress", func() {
			equalConfirmableEmailAddress, err := values.NewEmailAddress(emailAddressValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := confirmableEmailAddress.Equals(equalConfirmableEmailAddress)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given a different EmailAddress", func() {
			differentEmailAddress, err := values.NewEmailAddress("foo+different@bar.com")
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

/*** Tests for Modification methods ***/

func TestConfirmableEmailAddressConfirm(t *testing.T) {
	Convey("Given a new ConfirmableEmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)
		unconfirmedEmailAddress := emailAddress.ToConfirmable()

		Convey("When it is confirmed with the right ConfirmationHash", func() {
			confirmationHash := values.ReconstituteConfirmationHash(unconfirmedEmailAddress.ConfirmationHash())
			confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(emailAddress, confirmationHash)

			Convey("It should be confirmed", func() {
				So(err, ShouldBeNil)
				So(confirmedEmailAddress.IsConfirmed(), ShouldBeTrue)
				So(unconfirmedEmailAddress.IsConfirmed(), ShouldBeFalse)
			})
		})

		Convey("When it is confirmed with a wrong ConfirmationHash", func() {
			confirmationHash := values.ReconstituteConfirmationHash("invalid_confirmation_hash")
			confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(emailAddress, confirmationHash)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrNotEqual), ShouldBeTrue)
				So(confirmedEmailAddress, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Marshal/Unmarshal methods ***/

func TestConfirmableEmailAddressMarshalJSON(t *testing.T) {
	Convey("Given a ConfirmableEmailAddress", t, func() {
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
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
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
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
				So(err, ShouldNotBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
