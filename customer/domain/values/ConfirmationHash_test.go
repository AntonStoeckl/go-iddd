package values_test

import (
	"fmt"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateConfirmationHash(t *testing.T) {
	Convey("When a ConfirmationHash is generated", t, func() {
		confirmationHash := values.GenerateConfirmationHash("john@doe.com")

		Convey("It should succeed", func() {
			So(confirmationHash, ShouldHaveSameTypeAs, values.ConfirmationHash{})
			So(confirmationHash, ShouldNotBeZeroValue)
			So(confirmationHash.Hash(), ShouldNotBeBlank)

			_, err := values.BuildConfirmationHash(confirmationHash.Hash())
			So(err, ShouldBeNil)
		})
	})

	Convey("When many ConfirmationHashes are generated with similar input", t, func() {
		amount := 200
		hashes := make(map[string]int)

		for i := 0; i < amount; i++ {
			input := fmt.Sprintf("john+%d@doe.com", i)
			id := values.GenerateConfirmationHash(input)
			hashes[id.Hash()] = i
			time.Sleep(time.Nanosecond)
		}

		Convey("They should have unique values", func() {
			So(hashes, ShouldHaveLength, amount)
		})
	})
}

func TestBuildConfirmationHash(t *testing.T) {
	Convey("When a ConfirmationHash is built with valid input", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash, err := values.BuildConfirmationHash(confirmationHashValue)

		Convey("It should succeed", func() {
			So(err, ShouldBeNil)
			So(confirmationHash, ShouldNotBeNil)
			So(confirmationHash, ShouldHaveSameTypeAs, values.ConfirmationHash{})
			So(confirmationHash.Hash(), ShouldEqual, confirmationHashValue)
		})
	})

	Convey("When a ConfirmationHash is built with invalid input", t, func() {
		confirmationHashValue := ""
		confirmationHash, err := values.BuildConfirmationHash(confirmationHashValue)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
			So(confirmationHash, ShouldBeZeroValue)
		})
	})
}

func TestRebuildConfirmationHash(t *testing.T) {
	Convey("When a ConfirmationHash is rebuilt", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash := values.RebuildConfirmationHash(confirmationHashValue)

		Convey("It should succeed", func() {
			So(confirmationHash, ShouldHaveSameTypeAs, values.ConfirmationHash{})
			So(confirmationHash, ShouldNotBeZeroValue)
			So(confirmationHash.Hash(), ShouldEqual, confirmationHashValue)
		})
	})
}

func TestConfirmationHashEquals(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash, err := values.BuildConfirmationHash(confirmationHashValue)
		So(err, ShouldBeNil)

		Convey("And given an equal ConfirmationHash", func() {
			equalConfirmationHash, err := values.BuildConfirmationHash(confirmationHashValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := confirmationHash.Equals(equalConfirmationHash)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another different ConfirmationHash", func() {
			differentConfirmationHashValue := "different_hash"
			differentConfirmationHash, err := values.BuildConfirmationHash(differentConfirmationHashValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := confirmationHash.Equals(differentConfirmationHash)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}
