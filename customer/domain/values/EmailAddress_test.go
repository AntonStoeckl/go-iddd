package values_test

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

/*** Tests for Factory methods ***/

func TestNewEmailAddress(t *testing.T) {
	Convey("Given a valid emailAddress as input", t, func() {
		Convey("When a new EmailAddress is created", func() {
			emailAddress, err := values.NewEmailAddress("foo@bar.com")

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(emailAddress, ShouldNotBeNil)
				So(emailAddress, ShouldHaveSameTypeAs, (*values.EmailAddress)(nil))
			})
		})
	})

	Convey("Given an invalid emailAddress as input", t, func() {
		Convey("When a new EmailAddress is created", func() {
			emailAddress, err := values.NewEmailAddress("foo@bar.c")

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				So(emailAddress, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Getter methods ***/

func TestEmailAddressExposesExpectedValues(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			So(emailAddress.EmailAddress(), ShouldEqual, emailAddressValue)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestEmailAddressEquals(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)

		Convey("And given an equal EmailAddress", func() {
			equalEmailAddress, err := values.NewEmailAddress(emailAddressValue)
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
			differentEmailAddress, err := values.NewEmailAddress(differentEmailAddressValue)
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

func TestEmailAddressShouldEqual(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddressValue := "foo@bar.com"
		emailAddress, err := values.NewEmailAddress(emailAddressValue)
		So(err, ShouldBeNil)

		Convey("And given an equal EmailAddress", func() {
			equalEmailAddress, err := values.NewEmailAddress(emailAddressValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				err := emailAddress.ShouldEqual(equalEmailAddress)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("And given a different EmailAddress", func() {
			differentEmailAddress, err := values.NewEmailAddress("foo+different@bar.com")
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				err := emailAddress.ShouldEqual(differentEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, shared.ErrNotEqual), ShouldBeTrue)
				})
			})
		})
	})
}

/*** Tests for Conversion methods ***/

func TestEmailAddressToConfirmable(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		Convey("When it is converted to confirmable", func() {
			confirmableEmailAddress := emailAddress.ToConfirmable()

			Convey("It should be a ConfirmableEmailAddress", func() {
				So(err, ShouldBeNil)
				So(confirmableEmailAddress, ShouldNotBeNil)
				So(confirmableEmailAddress, ShouldHaveSameTypeAs, (*values.ConfirmableEmailAddress)(nil))
			})
		})
	})
}

/*** Tests for Marshal/Unmarshal methods ***/

func TestEmailAddressMarshalJSON(t *testing.T) {
	Convey("Given an EmailAddress", t, func() {
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		Convey("When it is marshaled to json", func() {
			data, err := emailAddress.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, `"`+emailAddress.EmailAddress()+`"`)
			})
		})
	})
}

func TestEmailAddressUnmarshalJSON(t *testing.T) {
	Convey("Given an EmailAddress marshaled to json", t, func() {
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		data, err := emailAddress.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &values.EmailAddress{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original EmailAddress", func() {
				So(err, ShouldBeNil)
				So(emailAddress, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to EmailAddress", func() {
			unmarshaled := &values.EmailAddress{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
