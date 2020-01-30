package application

import (
	"database/sql"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	EventStream(id values.CustomerID) (es.DomainEvents, error)
	Register(id values.CustomerID, recordedEvents es.DomainEvents, tx *sql.Tx) error
	Persist(id values.CustomerID, recordedEvents es.DomainEvents, tx *sql.Tx) error
	Delete(id values.CustomerID) error
}
