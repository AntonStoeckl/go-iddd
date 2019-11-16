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

/*** Tests for Factory methods ***/

func TestGenerateConfirmationHash(t *testing.T) {
	Convey("When a ConfirmationHash is generated", t, func() {
		confirmationHash := values.GenerateConfirmationHash("john@doe.com")

		Convey("It should succeed", func() {
			So(confirmationHash, ShouldNotBeNil)
			So(confirmationHash, ShouldHaveSameTypeAs, (*values.ConfirmationHash)(nil))
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

func TestRebuildConfirmationHash(t *testing.T) {
	Convey("Given that the supplied input is valid", t, func() {
		confirmationHashValue := "secret_hash"

		Convey("When a ConfirmationHash is rebuilt", func() {
			confirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(confirmationHash, ShouldNotBeNil)
				So(confirmationHash, ShouldHaveSameTypeAs, (*values.ConfirmationHash)(nil))
			})
		})
	})

	Convey("Given that the supplied input is not valid", t, func() {
		confirmationHashValue := ""

		Convey("When a ConfirmationHash is rebuilt", func() {
			confirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				So(confirmationHash, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Getter methods ***/

func TestConfirmationHashExposesExpectedValues(t *testing.T) {
	Convey("Given a generated ConfirmationHash", t, func() {
		confirmationHashInput := "foo@bar.com"
		confirmationHash := values.GenerateConfirmationHash(confirmationHashInput)

		Convey("It should expose a generated value", func() {
			So(confirmationHash.Hash(), ShouldNotBeBlank)
		})
	})

	Convey("Given a rebuilt ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)
		So(err, ShouldBeNil)

		Convey("It should expose the expected value", func() {
			So(confirmationHash.Hash(), ShouldEqual, confirmationHashValue)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestConfirmationHashShouldEqual(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)
		So(err, ShouldBeNil)

		Convey("And given an equal ConfirmationHash", func() {
			equalConfirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				Convey("They should be equal", func() {
					So(confirmationHash.Equals(equalConfirmationHash), ShouldBeTrue)
				})
			})
		})

		Convey("And given another different ConfirmationHash", func() {
			differentConfirmationHashValue := "different_hash"
			differentConfirmationHash, err := values.RebuildConfirmationHash(differentConfirmationHashValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				Convey("They should not be equal", func() {
					So(confirmationHash.Equals(differentConfirmationHash), ShouldBeFalse)
				})
			})
		})
	})
}

func TestConfirmationHashEquals(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)
		So(err, ShouldBeNil)

		Convey("And given an equal ConfirmationHash", func() {
			equalConfirmationHash, err := values.RebuildConfirmationHash(confirmationHashValue)
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
			differentConfirmationHash, err := values.RebuildConfirmationHash(differentConfirmationHashValue)
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

/*** Tests for Marshal/Unmarshal methods ***/

func TestConfirmationHashMarshalJSON(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHash := values.GenerateConfirmationHash("foo@bar.com")

		Convey("When it is marshaled to json", func() {
			data, err := confirmationHash.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, `"`+confirmationHash.Hash()+`"`)
			})
		})
	})
}

func TestConfirmationHashUnmarshalJSON(t *testing.T) {
	Convey("Given a ConfirmationHash marshaled to json", t, func() {
		confirmationHash := values.GenerateConfirmationHash("foo@bar.com")
		data, err := confirmationHash.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &values.ConfirmationHash{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original ConfirmationHash", func() {
				So(err, ShouldBeNil)
				So(confirmationHash, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to ConfirmationHash", func() {
			unmarshaled := &values.ConfirmationHash{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
