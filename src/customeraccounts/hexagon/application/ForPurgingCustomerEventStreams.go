package application

import "github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"

type ForPurgingCustomerEventStreams func(id value.CustomerID) error
