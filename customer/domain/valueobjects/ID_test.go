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
        })

        Convey("And then it should expose the expected id", func() {
            So(id.String(), ShouldNotBeBlank)
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
        from := "someId"
        id := ReconstituteID(from)

        Convey("Then it should reconstitute an ID", func() {
            So(id, ShouldNotBeNil)
        })

        Convey("And then it should expose the expected id", func() {
            So(id.String(), ShouldEqual, from)
        })
    })
}

func TestEqualsOnID(t *testing.T) {
    Convey("Given an Identifier of type ID", t, func() {
        from := "someId"
        id := ReconstituteID(from)

        Convey("And given a second ID with equal value", func() {
            equalId := ReconstituteID(from)

            Convey("When they are compared", func() {
                isEqual := id.Equals(equalId)

                Convey("Then they should be equal", func() {
                    So(isEqual, ShouldBeTrue)
                })
            })
        })

        Convey("And given a second Identifier with equal value but different type", func() {
            differentId := &dummyIdentifier{value: from}

            Convey("When they are compared", func() {
                isEqual := id.Equals(differentId)

                Convey("Then they should not be equal", func() {
                    So(isEqual, ShouldBeFalse)
                })
            })
        })

        Convey("And given a second ID with different value", func() {
            differentId := ReconstituteID("differentId")

            Convey("When they are compared", func() {
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
    return false
}
