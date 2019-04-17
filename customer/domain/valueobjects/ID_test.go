package valueobjects

import (
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

	Convey("When GenerateID is invoked multiple times", t, func() {
		ids := make(map[string]int)
		amount := 100000

		for i := 0; i < amount; i++ {
			id := GenerateID()
			ids[id.String()] = i
		}

		Convey("Then it should create unique ids", func() {
			So(ids, ShouldHaveLength, amount)
		})
	})
}

func TestReconstituteID(t *testing.T) {
	Convey("When ReconstituteID is invoked", t, func() {
		from := "ee58f146-a233-40ee-9a21-b57106676e72"
		id := ReconstituteID(from)

		Convey("Then it should reconstitute an ID", func() {
			So(id, ShouldNotBeNil)
		})

		Convey("And then it should expose the expected id", func() {
			So(id.String(), ShouldEqual, from)
		})
	})
}
