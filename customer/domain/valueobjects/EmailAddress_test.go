package valueobjects

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewEmailAddress(t *testing.T) {
	Convey("Given that the supplied emailAddress is valid", t, func() {
		emailAddressValue := "foo@bar.com"

		Convey("When NewEmailAddress is invoked", func() {
			emailAddress, err := NewEmailAddress(emailAddressValue)

			Convey("Then it should create an EmailAddress", func() {
				So(err, ShouldBeNil)
				So(emailAddress, ShouldImplement, (*EmailAddress)(nil))
			})

			Convey("And then it should expose the expected emailAddress", func() {
				So(emailAddress.String(), ShouldEqual, emailAddressValue)
			})
		})
	})

	Convey("Given that the supplied emailAddress is not valid", t, func() {
		emailAddressValue := "foo@bar.c"

		Convey("When NewEmailAddress is invoked", func() {
			emailAddress, err := NewEmailAddress(emailAddressValue)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "emailAddress - invalid input given")
				So(emailAddress, ShouldBeNil)
			})
		})
	})
}

func TestEqualsOnEmailAddress(t *testing.T) {
	Convey("Given that two EmailAddresses represent equal emailAddresses", t, func() {
		emailAddress := ReconstituteEmailAddress("foo@bar.com")
		equalEmailAddress := ReconstituteEmailAddress("foo@bar.com")

		Convey("When Equal is invoked", func() {
			isEqual := emailAddress.Equals(equalEmailAddress)

			Convey("Then they should be equal", func() {
				So(isEqual, ShouldBeTrue)
			})
		})
	})

	Convey("Given that tow EmailAddresses represent different emailAddresses", t, func() {
		emailAddress := ReconstituteEmailAddress("foo@bar.com")
		differentEmailAddress := ReconstituteEmailAddress("foo+different@bar.com")

		Convey("When Equal is invoked", func() {
			isEqual := emailAddress.Equals(differentEmailAddress)

			Convey("Then they should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})
	})
}
