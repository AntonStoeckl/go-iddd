package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForStartingCustomerEventStreams func(recordedEvents es.RecordedEvents, id value.CustomerID) error
