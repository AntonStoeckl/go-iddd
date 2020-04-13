package values_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPersonName_Equals(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		personName := values.RebuildPersonName("Lib", "Gallagher")

		Convey("When it is compared with an identical PersonName", func() {
			identicalPersonName := values.RebuildPersonName(personName.GivenName(), personName.FamilyName())
			isEqual := personName.Equals(identicalPersonName)

			Convey("Then it should be equal", func() {
				So(isEqual, ShouldBeTrue)
			})
		})

		Convey("When it is compared with another PersonName with different givenName", func() {
			differentPersonName := values.RebuildPersonName("Phillip", personName.FamilyName())
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})

		Convey("When it is compared with another PersonName with different familyName", func() {
			differentPersonName := values.RebuildPersonName(personName.GivenName(), "Jackson")
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})
	})
}
