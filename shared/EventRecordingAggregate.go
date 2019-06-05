package shared

type EventRecordingAggregate interface {
	Aggregate
	RecordsEvents
}
