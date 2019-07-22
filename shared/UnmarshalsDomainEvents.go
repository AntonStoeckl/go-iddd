package shared

type UnmarshalDomainEvent func(name string, payload []byte) (DomainEvent, error)
