package valueobjects_test

import (
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPersonName(t *testing.T) {
	givenName := "John"
	familyName := "Doe"

	Convey("Given that the supplied givenName and familyName are valid", t, func() {
		Convey("When NewPersonName is invoked", func() {
			personName, err := valueobjects.NewPersonName(givenName, familyName)

			Convey("Then it should create a PersonName", func() {
				So(err, ShouldBeNil)
				So(personName, ShouldNotBeNil)
				So(personName, ShouldHaveSameTypeAs, (*valueobjects.PersonName)(nil))
			})

			Convey("And then it should expose the expected values", func() {
				So(personName.GivenName(), ShouldEqual, givenName)
				So(personName.FamilyName(), ShouldEqual, familyName)
			})
		})
	})

	Convey("Given that the supplied givenName is not valid", t, func() {
		givenName = ""

		Convey("When NewPersonName is invoked", func() {
			emailAddress, err := valueobjects.NewPersonName(givenName, familyName)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "personName - empty input given for givenName")
				So(emailAddress, ShouldBeNil)
			})
		})
	})

	Convey("Given that the supplied familyName is not valid", t, func() {
		familyName = ""

		Convey("When NewPersonName is invoked", func() {
			emailAddress, err := valueobjects.NewPersonName(givenName, familyName)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError, "personName - empty input given for familyName")
				So(emailAddress, ShouldBeNil)
			})
		})
	})
}

func TestReconstitutePersonName(t *testing.T) {
	Convey("When ReconstitutePersonName invoked", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName := valueobjects.ReconstitutePersonName(givenName, familyName)

		Convey("Then it should reconstitute a PersonName", func() {
			So(personName, ShouldNotBeNil)
			So(personName, ShouldHaveSameTypeAs, (*valueobjects.PersonName)(nil))
		})

		Convey("And then it should expose the expected values", func() {
			So(personName.GivenName(), ShouldEqual, givenName)
			So(personName.FamilyName(), ShouldEqual, familyName)
		})
	})
}

func TestEqualsOnPersonName(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName := valueobjects.ReconstitutePersonName(givenName, familyName)

		Convey("And given another equal PersonName", func() {
			equalPersonName := valueobjects.ReconstitutePersonName(givenName, familyName)

			Convey("When Equals is invoked", func() {
				isEqual := personName.Equals(equalPersonName)

				Convey("Then they should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another PersonName with different givenName", func() {
			givenName = "Peter"
			differentPersonName := valueobjects.ReconstitutePersonName(givenName, familyName)

			Convey("When Equals is invoked", func() {
				isEqual := personName.Equals(differentPersonName)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given another PersonName with different familyName", func() {
			familyName = "Mueller"
			differentPersonName := valueobjects.ReconstitutePersonName(givenName, familyName)

			Convey("When Equals is invoked", func() {
				isEqual := personName.Equals(differentPersonName)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
