package command

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ForAssertingUniqueEmailAddresses interface {
	Assert(assertions customer.UniqueEmailAddressAssertions, tx *sql.Tx) error
	ClearFor(customerID values.CustomerID, tx *sql.Tx) error
}
