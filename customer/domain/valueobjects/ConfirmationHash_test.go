package valueobjects_test

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
	"testing"

	"golang.org/x/xerrors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateConfirmationHash(t *testing.T) {
	Convey("Given some input", t, func() {
		confirmationHashValue := "secret_hash"

		Convey("When a ConfirmationHash is generated", func() {
			confirmationHash := valueobjects.GenerateConfirmationHash(confirmationHashValue)

			Convey("It should succeed", func() {
				So(confirmationHash, ShouldNotBeNil)
				So(confirmationHash, ShouldHaveSameTypeAs, (*valueobjects.ConfirmationHash)(nil))
			})
		})
	})
}

func TestReconstituteConfirmationHash(t *testing.T) {
	Convey("Given some input", t, func() {
		confirmationHashValue := "secret_hash"

		Convey("When a ConfirmationHash is reconstituted", func() {
			confirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

			Convey("It should succeed", func() {
				So(confirmationHash, ShouldNotBeNil)
				So(confirmationHash, ShouldHaveSameTypeAs, (*valueobjects.ConfirmationHash)(nil))
			})
		})
	})
}

func TestConfirmationHashExposesExpectedValues(t *testing.T) {
	Convey("Given a generated ConfirmationHash", t, func() {
		confirmationHashInput := "foo@bar.com"
		confirmationHash := valueobjects.GenerateConfirmationHash(confirmationHashInput)

		Convey("It should expose a generated value", func() {
			So(confirmationHash.Hash(), ShouldNotBeBlank)
		})
	})

	Convey("Given a reconstituted ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

		Convey("It should expose the expected value", func() {
			So(confirmationHash.Hash(), ShouldEqual, confirmationHashValue)
		})
	})
}

func TestConfirmationHashShouldEqual(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHashValue := "secret_hash"
		confirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

		Convey("And given an equal ConfirmationHash", func() {
			equalConfirmationHash := valueobjects.ReconstituteConfirmationHash(confirmationHashValue)

			Convey("When they are compared", func() {
				err := confirmationHash.ShouldEqual(equalConfirmationHash)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("And given another different ConfirmationHash", func() {
			differentConfirmationHashValue := "different_hash"
			differentConfirmationHash := valueobjects.ReconstituteConfirmationHash(differentConfirmationHashValue)

			Convey("When they are compared", func() {
				err := confirmationHash.ShouldEqual(differentConfirmationHash)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(xerrors.Is(err, shared.ErrNotEqual), ShouldBeTrue)
				})
			})
		})
	})
}

func TestConfirmationHashMarshalJSON(t *testing.T) {
	Convey("Given a ConfirmationHash", t, func() {
		confirmationHash := valueobjects.GenerateConfirmationHash("foo@bar.com")

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
		confirmationHash := valueobjects.GenerateConfirmationHash("foo@bar.com")
		data, err := confirmationHash.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &valueobjects.ConfirmationHash{}
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
			unmarshaled := &valueobjects.ConfirmationHash{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
