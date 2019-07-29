package shared

type StartsEventStoreSessions interface {
	StartSession() (EventStoreSession, error)
}
