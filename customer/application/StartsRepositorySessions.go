package application

import "database/sql"

type StartsRepositorySessions interface {
	StartSession(tx *sql.Tx) PersistableCustomers
}
