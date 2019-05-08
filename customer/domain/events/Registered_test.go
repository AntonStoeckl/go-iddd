package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalJSONOnRegistered(t *testing.T) {
	Convey("Given a valid Registered event", t, func() {
		id := valueobjects.GenerateID()
		emailAddress, err := valueobjects.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := valueobjects.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		event := ItWasRegistered(id, confirmableEmailAddress, personName)
		So(event, ShouldNotBeNil)
		So(event, ShouldHaveSameTypeAs, (*Registered)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := json.Marshal(event)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})

			Convey("And when it is unmarshaled from json", func() {
				unmarshaledEvent := &Registered{}
				err := json.Unmarshal(data, unmarshaledEvent)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(unmarshaledEvent, ShouldNotBeNil)
					So(unmarshaledEvent, ShouldHaveSameTypeAs, (*Registered)(nil))
					So(event, ShouldResemble, unmarshaledEvent)
				})
			})
		})
	})
}
