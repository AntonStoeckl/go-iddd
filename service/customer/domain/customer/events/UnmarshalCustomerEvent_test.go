package events

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnmarshalCustomerEvent(t *testing.T) {
	Convey("When an unknown event is unmarshaled", t, func() {
		_, err := UnmarshalCustomerEvent("unknown", []byte{}, 1)
		Convey("Then it should fail", func() {
			So(errors.Is(err, lib.ErrUnmarshalingFailed), ShouldBeTrue)
		})
	})
}