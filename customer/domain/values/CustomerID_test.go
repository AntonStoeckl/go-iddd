package values_test

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"sync"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

/*** Tests for Factory methods ***/

func TestGenerateCustomerID(t *testing.T) {
	Convey("When an CustomerID is generated", t, func() {
		customerID := values.GenerateCustomerID()

		Convey("It should succeed", func() {
			So(customerID, ShouldNotBeNil)
			So(customerID, ShouldHaveSameTypeAs, (*values.CustomerID)(nil))
		})
	})

	Convey("When many IDs are generated", t, func() {
		group := sync.WaitGroup{}
		mutex := sync.Mutex{}
		customerIDs := make(map[string]int)
		amountPerRoutine := 100
		numRoutines := 500
		totalAmount := 0

		for i := 0; i < numRoutines; i++ {
			go generateManyCustomerIDs(customerIDs, &group, &mutex, amountPerRoutine)
			group.Add(1)
			totalAmount += amountPerRoutine
		}

		group.Wait()

		Convey("They should have unique values", func() {
			So(customerIDs, ShouldHaveLength, totalAmount)
		})
	})
}

func generateManyCustomerIDs(ids map[string]int, group *sync.WaitGroup, mutex *sync.Mutex, amountPerRoutine int) {
	generatedIDs := make(map[string]int)

	for i := 0; i < amountPerRoutine; i++ {
		id := values.GenerateCustomerID()
		generatedIDs[id.String()] = i
	}

	mutex.Lock()
	for key, value := range generatedIDs {
		ids[key] = value
	}
	mutex.Unlock()

	group.Done()
}

func TestRebuildCustomerID(t *testing.T) {
	Convey("Given that the supplied id is valid", t, func() {
		idValue := "b5f1a1b1-5d03-4e08-8365-259791228be3"

		Convey("When an CustomerID is rebuilt", func() {
			customerID, err := values.RebuildCustomerID(idValue)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(customerID, ShouldNotBeNil)
				So(customerID, ShouldHaveSameTypeAs, (*values.CustomerID)(nil))
			})
		})
	})

	Convey("Given that the supplied id is not valid", t, func() {
		idValue := ""

		Convey("When an CustomerID is rebuilt", func() {
			customerID, err := values.RebuildCustomerID(idValue)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				So(customerID, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Getter methods ***/

func TestCustomerIDExposesExpectedValues(t *testing.T) {
	Convey("Given a generated CustomerID", t, func() {
		customerID := values.GenerateCustomerID()

		Convey("It should expose a generated value", func() {
			So(customerID.String(), ShouldNotBeBlank)
		})
	})

	Convey("Given a rebuilt CustomerID", t, func() {
		idValue := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		customerID, err := values.RebuildCustomerID(idValue)
		So(err, ShouldBeNil)

		Convey("It should expose the expected value", func() {
			So(customerID.String(), ShouldEqual, idValue)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestEqualsOnCustomerID(t *testing.T) {
	Convey("Given an Identifier of type CustomerID", t, func() {
		idValue := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		customerID, err := values.RebuildCustomerID(idValue)
		So(err, ShouldBeNil)

		Convey("And given an equal CustomerID", func() {
			equalId, err := values.RebuildCustomerID(idValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := customerID.Equals(equalId)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given an CustomerID with equal value but different type", func() {
			differentId := &dummyIdentifier{value: idValue}

			Convey("When they are compared", func() {
				isEqual := customerID.Equals(differentId)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given an CustomerID with equal type but different value", func() {
			differentIdValue := "5b6e0bc9-aa69-4dd9-be1c-d54bee80f565"
			differentId, err := values.RebuildCustomerID(differentIdValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := customerID.Equals(differentId)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}

/*** Tests for Marshal/Unmarshal methods ***/

func TestCustomerIDMarshalJSON(t *testing.T) {
	Convey("Given an CustomerID", t, func() {
		customerID := values.GenerateCustomerID()

		Convey("When it is marshaled to json", func() {
			data, err := customerID.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, `"`+customerID.String()+`"`)
			})
		})
	})
}

func TestCustomerIDUnmarshalJSON(t *testing.T) {
	Convey("Given an CustomerID marshaled to json", t, func() {
		customerID := values.GenerateCustomerID()
		data, err := customerID.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &values.CustomerID{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original CustomerID", func() {
				So(err, ShouldBeNil)
				So(customerID, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to CustomerID", func() {
			unmarshaled := &values.CustomerID{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}

/*** Test helper type implementing shared.IdentifiesAggregates ***/

type dummyIdentifier struct {
	value string
}

func (idenfifier *dummyIdentifier) String() string {
	return idenfifier.value
}

func (idenfifier *dummyIdentifier) Equals(other shared.IdentifiesAggregates) bool {
	// this method will never be invoked, because we don't test this dummy
	return false
}
