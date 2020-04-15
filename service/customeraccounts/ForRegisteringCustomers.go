package customeraccounts

import "github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"

type ForRegisteringCustomers func(emailAddress, givenName, familyName string) (value.CustomerID, error)
