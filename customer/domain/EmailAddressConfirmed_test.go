package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalJSONOnEmailAddressConfirmed(t *testing.T) {
	Convey("Given a valid EmailAddressConfirmed event", t, func() {
		id := valueobjects.GenerateID()
		emailAddress, err := valueobjects.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)

		event := EmailAddressWasConfirmed(id, emailAddress)
		So(event, ShouldNotBeNil)
		So(event, ShouldHaveSameTypeAs, (*EmailAddressConfirmed)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := json.Marshal(event)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})

			Convey("And when it is unmarshaled from json", func() {
				unmarshaledEvent := &EmailAddressConfirmed{}
				err := json.Unmarshal(data, unmarshaledEvent)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(unmarshaledEvent, ShouldNotBeNil)
					So(unmarshaledEvent, ShouldHaveSameTypeAs, (*EmailAddressConfirmed)(nil))
					So(event, ShouldResemble, unmarshaledEvent)
				})
			})
		})
	})
}
