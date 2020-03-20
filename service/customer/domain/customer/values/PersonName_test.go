package values_test

import (
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildPersonName(t *testing.T) {
	Convey("When a PersonName is built from valid input", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName, err := values.BuildPersonName(givenName, familyName)

		Convey("It should succeed", func() {
			So(err, ShouldBeNil)
			So(personName, ShouldHaveSameTypeAs, values.PersonName{})
			So(personName, ShouldNotBeZeroValue)
			So(personName.GivenName(), ShouldEqual, givenName)
			So(personName.FamilyName(), ShouldEqual, familyName)
		})
	})

	Convey("When a PersonName is built with invalid givenName", t, func() {
		givenName := ""
		familyName := "Doe"
		personName, err := values.BuildPersonName(givenName, familyName)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
			So(personName, ShouldBeZeroValue)
		})
	})

	Convey("When a PersonName is built with invalid familyName", t, func() {
		givenName := "John"
		familyName := ""
		personName, err := values.BuildPersonName(givenName, familyName)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
			So(personName, ShouldBeZeroValue)
		})
	})
}

func TestRebuildPersonName(t *testing.T) {
	Convey("When a PersonName is built from valid input", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName := values.RebuildPersonName(givenName, familyName)

		Convey("It should succeed", func() {
			So(personName, ShouldHaveSameTypeAs, values.PersonName{})
			So(personName, ShouldNotBeZeroValue)
			So(personName.GivenName(), ShouldEqual, givenName)
			So(personName.FamilyName(), ShouldEqual, familyName)
		})
	})
}

func TestPersonNameEquals(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName, err := values.BuildPersonName(givenName, familyName)
		So(err, ShouldBeNil)

		Convey("And given an equal PersonName", func() {
			equalPersonName, err := values.BuildPersonName(givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := personName.Equals(equalPersonName)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given a PersonName with different givenName", func() {
			givenName = "Peter"
			differentPersonName, err := values.BuildPersonName(givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := personName.Equals(differentPersonName)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given a PersonName with different familyName", func() {
			familyName = "Mueller"
			differentPersonName, err := values.BuildPersonName(givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := personName.Equals(differentPersonName)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
