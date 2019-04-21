package valueobjects

import (
    "testing"

    . "github.com/smartystreets/goconvey/convey"
)

func TestNewConfirmableEmailAddress(t *testing.T) {
    Convey("Given that the supplied emailAddress is valid", t, func() {
        emailAddressValue := "foo@bar.com"

        Convey("When NewConfirmableEmailAddress is invoked", func() {
            confirmableEmailAddress, err := NewConfirmableEmailAddress(emailAddressValue)

            Convey("Then it should create a ConfirmableEmailAddress", func() {
                So(err, ShouldBeNil)
                So(confirmableEmailAddress, ShouldImplement, (*ConfirmableEmailAddress)(nil))
            })

            Convey("And then it should expose the expected emailAddress", func() {
                So(confirmableEmailAddress.String(), ShouldEqual, emailAddressValue)
            })

            Convey("And then it should not be confirmed", func() {
                So(confirmableEmailAddress.IsConfirmed(), ShouldBeFalse)
            })
        })
    })

    Convey("Given that the supplied emailAddress is not valid", t, func() {
        emailAddressValue := "foo@bar.c"

        Convey("When NewConfirmableEmailAddress is invoked", func() {
            confirmableEmailAddress, err := NewConfirmableEmailAddress(emailAddressValue)

            Convey("Then it should fail", func() {
                So(err, ShouldBeError, "emailAddress - invalid input given")
                So(confirmableEmailAddress, ShouldBeNil)
            })
        })
    })
}

func TestConfirmOnConfirmableEmailAddress(t *testing.T) {
    Convey("Given an unconfirmed EmailAddress", t, func() {
        unconfirmedEmailAddress := ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")

        Convey("When Confirm is invoked with a matching ConfirmationHash", func() {
            suppliedConfirmationHash := ReconstituteConfirmationHash("secret_hash")
            confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(suppliedConfirmationHash)

            Convey("Then it should succeed", func() {
                So(err, ShouldBeNil)
                So(confirmedEmailAddress.Equals(unconfirmedEmailAddress), ShouldBeTrue)
            })

            Convey("And then it should be confirmed", func() {
                So(confirmedEmailAddress.IsConfirmed(), ShouldBeTrue)
            })

            Convey("And then the original should be unchanged", func() {
                So(unconfirmedEmailAddress.IsConfirmed(), ShouldBeFalse)
            })

        })

        Convey("When Confirm is invoked with a different ConfirmationHash", func() {
            suppliedConfirmationHash := ReconstituteConfirmationHash("different_hash")
            confirmedEmailAddress, err := unconfirmedEmailAddress.Confirm(suppliedConfirmationHash)

            Convey("Then it should fail", func() {
                So(err, ShouldBeError, "confirmationHash - is not equal")
                So(confirmedEmailAddress, ShouldBeNil)
            })
        })
    })
}

func TestEqualsOnConfirmableEmailAddress(t *testing.T) {
    Convey("Given that two ConfirmableEmailAddresses represent equal emailAddresses", t, func() {
        emailAddress := ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
        equalEmailAddress := ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")

        Convey("When Equal is invoked", func() {
            isEqual := emailAddress.Equals(equalEmailAddress)

            Convey("Then they should be equal", func() {
                So(isEqual, ShouldBeTrue)
            })
        })
    })
    Convey("Given that tow ConfirmableEmailAddresses represent different emailAddresses", t, func() {
        emailAddress := ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
        differentEmailAddress := ReconstituteConfirmableEmailAddress("foo+different@bar.com", "secret_hash")

        Convey("When Equal is invoked", func() {
            isEqual := emailAddress.Equals(differentEmailAddress)

            Convey("Then they should not be equal", func() {
                So(isEqual, ShouldBeFalse)
            })
        })
    })
}
