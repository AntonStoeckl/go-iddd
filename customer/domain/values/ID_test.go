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

func TestGenerateID(t *testing.T) {
	Convey("When an ID is generated", t, func() {
		id := values.GenerateID()

		Convey("It should succeed", func() {
			So(id, ShouldNotBeNil)
			So(id, ShouldHaveSameTypeAs, (*values.ID)(nil))
		})
	})

	Convey("When many IDs are generated", t, func() {
		group := sync.WaitGroup{}
		mutex := sync.Mutex{}
		ids := make(map[string]int)
		amountPerRoutine := 100
		numRoutines := 500
		totalAmount := 0

		for i := 0; i < numRoutines; i++ {
			go generateManyIDs(ids, &group, &mutex, amountPerRoutine)
			group.Add(1)
			totalAmount += amountPerRoutine
		}

		group.Wait()

		Convey("They should have unique values", func() {
			So(ids, ShouldHaveLength, totalAmount)
		})
	})
}

func generateManyIDs(ids map[string]int, group *sync.WaitGroup, mutex *sync.Mutex, amountPerRoutine int) {
	generatedIDs := make(map[string]int)

	for i := 0; i < amountPerRoutine; i++ {
		id := values.GenerateID()
		generatedIDs[id.String()] = i
	}

	mutex.Lock()
	for key, value := range generatedIDs {
		ids[key] = value
	}
	mutex.Unlock()

	group.Done()
}

func TestRebuildID(t *testing.T) {
	Convey("Given that the supplied id is valid", t, func() {
		idValue := "b5f1a1b1-5d03-4e08-8365-259791228be3"

		Convey("When an ID is rebuilt", func() {
			id, err := values.RebuildID(idValue)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(id, ShouldNotBeNil)
				So(id, ShouldHaveSameTypeAs, (*values.ID)(nil))
			})
		})
	})

	Convey("Given that the supplied id is not valid", t, func() {
		idValue := ""

		Convey("When an ID is rebuilt", func() {
			id, err := values.RebuildID(idValue)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
				So(id, ShouldBeNil)
			})
		})
	})
}

/*** Tests for Getter methods ***/

func TestIDExposesExpectedValues(t *testing.T) {
	Convey("Given a generated ID", t, func() {
		id := values.GenerateID()

		Convey("It should expose a generated value", func() {
			So(id.String(), ShouldNotBeBlank)
		})
	})

	Convey("Given a rebuilt ID", t, func() {
		idValue := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		id, err := values.RebuildID(idValue)
		So(err, ShouldBeNil)

		Convey("It should expose the expected value", func() {
			So(id.String(), ShouldEqual, idValue)
		})
	})
}

/*** Tests for Comparison methods ***/

func TestEqualsOnID(t *testing.T) {
	Convey("Given an Identifier of type ID", t, func() {
		idValue := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		id, err := values.RebuildID(idValue)
		So(err, ShouldBeNil)

		Convey("And given an equal ID", func() {
			equalId, err := values.RebuildID(idValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := id.Equals(equalId)

				Convey("They should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given an ID with equal value but different type", func() {
			differentId := &dummyIdentifier{value: idValue}

			Convey("When they are compared", func() {
				isEqual := id.Equals(differentId)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given an ID with equal type but different value", func() {
			differentIdValue := "5b6e0bc9-aa69-4dd9-be1c-d54bee80f565"
			differentId, err := values.RebuildID(differentIdValue)
			So(err, ShouldBeNil)

			Convey("When they are compared", func() {
				isEqual := id.Equals(differentId)

				Convey("They should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}

/*** Tests for Marshal/Unmarshal methods ***/

func TestIDMarshalJSON(t *testing.T) {
	Convey("Given an ID", t, func() {
		id := values.GenerateID()

		Convey("When it is marshaled to json", func() {
			data, err := id.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, `"`+id.String()+`"`)
			})
		})
	})
}

func TestIDUnmarshalJSON(t *testing.T) {
	Convey("Given an ID marshaled to json", t, func() {
		id := values.GenerateID()
		data, err := id.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &values.ID{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original ID", func() {
				So(err, ShouldBeNil)
				So(id, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to ID", func() {
			unmarshaled := &values.ID{}
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
