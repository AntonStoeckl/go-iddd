package eventstore

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ForCheckingUniqueEmailAddresses interface {
	AddUniqueEmailAddress(emailAddress values.EmailAddress, tx *sql.Tx) error
	ReplaceUniqueEmailAddress(previousEmailAddress values.EmailAddress, newEmailAddress values.EmailAddress, tx *sql.Tx) error
	RemoveUniqueEmailAddress(newEmailAddress values.EmailAddress, tx *sql.Tx) error
}
