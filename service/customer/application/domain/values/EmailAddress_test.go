package values_test

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildEmailAddress(t *testing.T) {
	Convey("When an EmailAddress is built from valid input", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.BuildEmailAddress(emailAddressValue)

		Convey("It should succeed", func() {
			So(err, ShouldBeNil)
			So(emailAddress, ShouldHaveSameTypeAs, values.EmailAddress{})
			So(emailAddress, ShouldNotBeZeroValue)
			So(emailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})

	Convey("When an EmailAddress is built from invalid input", t, func() {
		emailAddress, err := values.BuildEmailAddress("foo@bar.c")

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
			So(emailAddress, ShouldBeZeroValue)
		})
	})
}

func TestRebuildEmailAddress(t *testing.T) {
	Convey("When an EmailAddress is rebuilt", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress := values.RebuildEmailAddress(emailAddressValue)

		Convey("It should succeed", func() {
			So(emailAddress, ShouldHaveSameTypeAs, values.EmailAddress{})
			So(emailAddress, ShouldNotBeZeroValue)
			So(emailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})
}

func TestEmailAddressEquals(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.BuildEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)

		Convey("And given an equal EmailAddress", func() {
			equalEmailAddress, err := values.BuildEmailAddress(emailAddressValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := emailAddress.Equals(equalEmailAddress)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given a different EmailAddress", func() {
			differentEmailAddressValue := "foo+different@bar.com"
			differentEmailAddress, err := values.BuildEmailAddress(differentEmailAddressValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := emailAddress.Equals(differentEmailAddress)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
