package identity

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

func ComparePasswords(stream es.EventStream, plainPassword value.PlainPassword) error {
	return nil
}
