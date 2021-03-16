package serialization

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey" // nolint:stylecheck,golint
)

/***** a mock event to test marshaling unknown event *****/

type SomeEvent struct{}

func (event SomeEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta("SomeEvent", "never", "someID", "someID", 1)
}

func (event SomeEvent) IsFailureEvent() bool {
	return false
}

func (event SomeEvent) FailureReason() error {
	return nil
}

func AssertEventMetaResembles(originalEvent, unmarshaledEvent es.DomainEvent) {
	So(unmarshaledEvent.Meta().EventName(), ShouldEqual, originalEvent.Meta().EventName())
	So(unmarshaledEvent.Meta().OccurredAt(), ShouldEqual, originalEvent.Meta().OccurredAt())
	So(unmarshaledEvent.Meta().CausationID(), ShouldEqual, originalEvent.Meta().CausationID())
	So(unmarshaledEvent.Meta().StreamVersion(), ShouldEqual, originalEvent.Meta().StreamVersion())
	So(unmarshaledEvent.IsFailureEvent(), ShouldEqual, originalEvent.IsFailureEvent())
	So(unmarshaledEvent.FailureReason(), ShouldBeError)
	So(unmarshaledEvent.FailureReason().Error(), ShouldEqual, originalEvent.FailureReason().Error())
	So(errors.Is(originalEvent.FailureReason(), shared.ErrDomainConstraintsViolation), ShouldBeTrue)
	So(errors.Is(unmarshaledEvent.FailureReason(), shared.ErrDomainConstraintsViolation), ShouldBeTrue)
}
