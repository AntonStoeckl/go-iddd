package application

import "database/sql"

type StartsCustomersSession interface {
	StartSession(tx *sql.Tx) Customers
}
