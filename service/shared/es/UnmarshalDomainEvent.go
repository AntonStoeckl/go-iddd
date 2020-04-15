package es

type UnmarshalDomainEvent func(name string, payload []byte, streamVersion uint) (DomainEvent, error)
