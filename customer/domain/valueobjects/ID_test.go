package valueobjects

import (
	"go-iddd/shared"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateID(t *testing.T) {
	Convey("When GenerateID is invoked", t, func() {
		id := GenerateID()

		Convey("Then it should generate an ID", func() {
			So(id, ShouldNotBeNil)
			So(id, ShouldImplement, (*ID)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(id.ID(), ShouldNotBeBlank)
			So(id.String(), ShouldNotBeBlank)
			So(id.ID(), ShouldEqual, id.String())
		})
	})

	Convey("When GenerateID is invoked many times", t, func() {
		group := sync.WaitGroup{}
		mutex := sync.Mutex{}
		ids := make(map[string]int)
		amountPerRoutine := 100
		numRoutines := 1000
		totalAmount := 0

		for i := 0; i < numRoutines; i++ {
			go generateIDs(ids, &group, &mutex, amountPerRoutine)
			group.Add(1)
			totalAmount += amountPerRoutine
		}

		group.Wait()

		Convey("Then it should create unique ids", func() {
			So(ids, ShouldHaveLength, totalAmount)
		})
	})
}

func generateIDs(ids map[string]int, group *sync.WaitGroup, mutex *sync.Mutex, amountPerRoutine int) {
	generatedIDs := make(map[string]int)

	for i := 0; i < amountPerRoutine; i++ {
		id := GenerateID()
		generatedIDs[id.String()] = i
	}

	mutex.Lock()
	for key, value := range generatedIDs {
		ids[key] = value
	}
	mutex.Unlock()

	group.Done()
}

func TestReconstituteID(t *testing.T) {
	Convey("When ReconstituteID is invoked", t, func() {
		idValue := "b5f1a1b1-5d03-4e08-8365-259791228be3"
		id := ReconstituteID(idValue)

		Convey("Then it should reconstitute an ID", func() {
			So(id, ShouldNotBeNil)
			So(id, ShouldImplement, (*ID)(nil))
		})

		Convey("And then it should expose the expected value", func() {
			So(id.String(), ShouldEqual, idValue)
		})
	})
}

func TestEqualsOnID(t *testing.T) {
	Convey("Given an Identifier of type ID", t, func() {
		idValue := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		id := ReconstituteID(idValue)

		Convey("And given another ID with equal value", func() {
			equalId := ReconstituteID(idValue)

			Convey("When Equals is invoked", func() {
				isEqual := id.Equals(equalId)

				Convey("Then they should be equal", func() {
					So(isEqual, ShouldBeTrue)
				})
			})
		})

		Convey("And given another ID with equal value but different type", func() {
			differentId := &dummyIdentifier{value: idValue}

			Convey("When Equals is invoked", func() {
				isEqual := id.Equals(differentId)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})

		Convey("And given another ID with different value", func() {
			differentIdValue := "5b6e0bc9-aa69-4dd9-be1c-d54bee80f565"
			differentId := ReconstituteID(differentIdValue)

			Convey("When Equals is invoked", func() {
				isEqual := id.Equals(differentId)

				Convey("Then they should not be equal", func() {
					So(isEqual, ShouldBeFalse)
				})
			})
		})
	})
}

type dummyIdentifier struct {
	value string
}

func (idenfifier *dummyIdentifier) String() string {
	return idenfifier.value
}

func (idenfifier *dummyIdentifier) Equals(other shared.AggregateIdentifier) bool {
	// this method will never be invoked, as we don't test the dummy
	return false
}
