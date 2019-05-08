package valueobjects_test

import (
	"fmt"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

/*** Tests for Factory methods ***/

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
			personName, err := valueobjects.NewPersonName(givenName, familyName)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrInvalidInput), ShouldBeTrue)
				So(personName, ShouldBeNil)
			})
		})
	})

	Convey("Given that the supplied familyName is not valid", t, func() {
		familyName = ""

		Convey("When NewPersonName is invoked", func() {
			personName, err := valueobjects.NewPersonName(givenName, familyName)

			Convey("Then it should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrInvalidInput), ShouldBeTrue)
				So(personName, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Getter methods ***/

func TestPersonNameExposesExpectedValues(t *testing.T) {
	Convey("Given an PersonName", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName, err := valueobjects.NewPersonName(givenName, familyName)
		So(err, ShouldBeNil)

		Convey("It should expose the expected GivenName", func() {
			So(personName.GivenName(), ShouldEqual, givenName)
		})

		Convey("It should expose the expected FamilyName", func() {
			So(personName.FamilyName(), ShouldEqual, familyName)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestPersonNameEquals(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		givenName := "John"
		familyName := "Doe"
		personName, err := valueobjects.NewPersonName(givenName, familyName)
		So(err, ShouldBeNil)

		Convey("And given an equal PersonName", func() {
			equalPersonName, err := valueobjects.NewPersonName(givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := personName.Equals(equalPersonName)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another PersonName with different givenName", func() {
			givenName = "Peter"
			differentPersonName, err := valueobjects.NewPersonName(givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := personName.Equals(differentPersonName)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given another PersonName with different familyName", func() {
			familyName = "Mueller"
			differentPersonName, err := valueobjects.NewPersonName(givenName, familyName)
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

/*** Tests for Marshal/Unmarshal methods ***/

func TestPersonNameMarshalJSON(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		personName, err := valueobjects.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		Convey("When it is marshaled to json", func() {
			data, err := personName.MarshalJSON()

			expectedJSON := fmt.Sprintf(
				`{"givenName":"%s","familyName":"%s"}`,
				personName.GivenName(),
				personName.FamilyName(),
			)

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, expectedJSON)
			})
		})
	})
}

func TestPersonNameUnmarshalJSON(t *testing.T) {
	Convey("Given a PersonName marshaled to json", t, func() {
		personName, err := valueobjects.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)
		data, err := personName.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &valueobjects.PersonName{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original PersonName", func() {
				So(err, ShouldBeNil)
				So(personName, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to PersonName", func() {
			unmarshaled := &valueobjects.PersonName{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
