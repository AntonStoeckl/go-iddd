package application

import (
	"database/sql"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
)

type ForStoringCustomers interface {
	EventStream(id values.CustomerID) (lib.DomainEvents, error)
	Register(id values.CustomerID, recordedEvents lib.DomainEvents, tx *sql.Tx) error
	Persist(id values.CustomerID, recordedEvents lib.DomainEvents, tx *sql.Tx) error
	Delete(id values.CustomerID) error
}
