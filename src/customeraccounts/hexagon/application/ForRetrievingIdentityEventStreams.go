package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ForRetrievingIdentityEventStreams func(id value.IdentityID) (es.EventStream, error)

type ForStoringIdentityEventStreams interface {
	RetrieveIdentityEventStream(id value.IdentityID) (es.EventStream, error)
}
