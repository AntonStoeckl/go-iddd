package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
)

type ForStartingCustomerEventStreams func(customerRegistered domain.CustomerRegistered) error
