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
		emailAddress := valueobjects.ReconstituteConfirmableEmailAddress("foo@bar.com", "secret_hash")
		personName := valueobjects.ReconstitutePersonName("John", "Doe")

		event := ItWasRegistered(id, emailAddress, personName)
		So(event, ShouldNotBeNil)
		So(event, ShouldImplement, (*Registered)(nil))

		Convey("When it is marshaled to json", func() {
			data, err := json.Marshal(event)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})

			Convey("And when it is unmarshaled from json", func() {
				unmarshaledEvent, err := UnmarshalRegisteredFromJSON(data)

				Convey("Then it should succeed", func() {
					So(err, ShouldBeNil)
					So(unmarshaledEvent, ShouldNotBeNil)
					So(unmarshaledEvent, ShouldImplement, (*Registered)(nil))
					So(event, ShouldResemble, unmarshaledEvent)
				})
			})
		})
	})
}
