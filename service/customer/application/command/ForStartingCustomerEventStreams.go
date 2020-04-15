package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForStartingCustomerEventStreams func(recordedEvents es.RecordedEvents, id values.CustomerID) error
