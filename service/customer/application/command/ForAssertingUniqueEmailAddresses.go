package command

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForAssertingUniqueEmailAddresses interface {
	Assert(recordedEvents es.DomainEvents, tx *sql.Tx) error
	ClearFor(customerID values.CustomerID, tx *sql.Tx) error
}
