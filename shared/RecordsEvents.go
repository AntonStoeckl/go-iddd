package shared

type RecordsEvents interface {
	RecordedEvents(purge bool) DomainEvents
}
