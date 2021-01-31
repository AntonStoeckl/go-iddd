package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
)

type ForStartingCustomerEventStreams func(customerRegistered domain.CustomerRegistered) error
