package hexagon

import "github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"

type ForRegisteringCustomers func(emailAddress, givenName, familyName string) (value.CustomerID, error)
