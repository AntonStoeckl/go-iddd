package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

type ForStartingCustomerEventStreams func(recordedEvents es.RecordedEvents, id value.CustomerID) error
