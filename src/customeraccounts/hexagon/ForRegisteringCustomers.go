package hexagon

import "github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"

type ForRegisteringCustomers func(customerIDValue value.CustomerID, emailAddress, givenName, familyName string) error
