package value_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPersonName_Equals(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		personName := value.RebuildPersonName("Lib", "Gallagher")

		Convey("When it is compared with an identical PersonName", func() {
			identicalPersonName := value.RebuildPersonName(personName.GivenName(), personName.FamilyName())
			isEqual := personName.Equals(identicalPersonName)

			Convey("Then it should be equal", func() {
				So(isEqual, ShouldBeTrue)
			})
		})

		Convey("When it is compared with another PersonName with different givenName", func() {
			differentPersonName := value.RebuildPersonName("Phillip", personName.FamilyName())
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})

		Convey("When it is compared with another PersonName with different familyName", func() {
			differentPersonName := value.RebuildPersonName(personName.GivenName(), "Jackson")
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})
	})
}
