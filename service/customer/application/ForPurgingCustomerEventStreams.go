package application

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"

type ForPurgingCustomerEventStreams func(id value.CustomerID) error
