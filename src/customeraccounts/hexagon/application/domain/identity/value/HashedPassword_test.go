package value_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHashedPassword(t *testing.T) {
	Convey("When HashedPasswordFromPlainPassword with valid PlainPassword", t, func() {
		plainPW, err := value.BuildPlainPassword("superSecurePW123")
		So(err, ShouldBeNil)

		hashedPW, err := value.HashedPasswordFromPlainPassword(plainPW)

		Convey("Then it should succeed", func() {
			So(err, ShouldBeNil)

			Convey("And the HashedPassword should equal the PlainPassword", func() {
				So(hashedPW.String(), ShouldNotBeEmpty)
				So(hashedPW.CompareWith(plainPW), ShouldBeTrue)

				Convey("And when RebuildHashedPassword", func() {
					hashedPW2 := value.RebuildHashedPassword(hashedPW.String())

					Convey("Then it should equal the PlainPassword", func() {
						So(hashedPW2.CompareWith(plainPW), ShouldBeTrue)
					})
				})
			})

			Convey("And the HashedPassword should not equal a different PlainPassword", func() {
				plainPW2, err := value.BuildPlainPassword("superSecurePW12")
				So(err, ShouldBeNil)

				So(hashedPW.CompareWith(plainPW2), ShouldBeFalse)
			})
		})
	})
}
