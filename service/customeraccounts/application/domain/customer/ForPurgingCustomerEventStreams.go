package customer

import "github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"

type ForPurgingCustomerEventStreams func(id value.CustomerID) error
