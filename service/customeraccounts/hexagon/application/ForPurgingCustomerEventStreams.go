package application

import "github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"

type ForPurgingCustomerEventStreams func(id value.CustomerID) error
