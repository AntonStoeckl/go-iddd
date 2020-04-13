package command

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

type ForPurgingCustomerEventStreams func(id values.CustomerID) error
