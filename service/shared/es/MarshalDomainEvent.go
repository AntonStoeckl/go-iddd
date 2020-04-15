package es

type MarshalDomainEvent func(event DomainEvent) ([]byte, error)
