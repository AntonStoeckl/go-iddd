package valueobjects_test

import (
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewEmailAddress(t *testing.T) {
	Convey("Given that the supplied emailAddress is valid", t, func() {
		emailAddressValue := "foo@bar.com"

		Convey("When NewEmailAddress is invoked", func() {
			emailAddress, err := valueobjects.NewEmailAddress(emailAddressValue)

			Convey("Then it should create an EmailAddress", func() {
				So(err, ShouldBeNil)
				So(emailAddress, ShouldNotBeNil)
				So(emailAddress, ShouldImplement, (*valueobjects.EmailAddress)(nil))
			})

			Convey("And then it should expose the expected value", func() {
				So(emailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
			})
		})
	})

	Convey("Given that the supplied emailAddress is not valid", t, func() {
		emailAddressValue := "foo@bar.c"

		Convey("When NewEmailAddress is invoked", func() {
			emailAddress, err := valueobjects.NewEmailAddress(emailAddressValue)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "emailAddress - invalid input given")
				So(emailAddress, ShouldBeNil)
			})
		})
	})
}

func TestReconstituteEmailAddress(t *testing.T) {
	Convey("When ReconstituteEmailAddress invoked", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress := valueobjects.ReconstituteEmailAddress(emailAddressValue)

		Convey("Then it should reconstitute an EmailAddress", func() {
			So(emailAddress, ShouldNotBeNil)
			So(emailAddress, ShouldImplement, (*valueobjects.EmailAddress)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(emailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})
}

func TestEqualsOnEmailAddress(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress := valueobjects.ReconstituteEmailAddress(emailAddressValue)

		Convey("And given another equal EmailAddress", func() {
			equalEmailAddress := valueobjects.ReconstituteEmailAddress(emailAddressValue)

			Convey("When Equals is invoked", func() {
				isEqual := emailAddress.Equals(equalEmailAddress)

				Convey("Then they should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another different EmailAddress", func() {
			differentEmailAddressValue := "foo+different@bar.com"
			differentEmailAddress := valueobjects.ReconstituteEmailAddress(differentEmailAddressValue)

			Convey("When Equals is invoked", func() {
				isEqual := emailAddress.Equals(differentEmailAddress)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
