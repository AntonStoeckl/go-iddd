package es

type DomainEvent interface {
	Meta() EventMeta
	IndicatesAnError() (bool, string)
}
