package application

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"

type ForRegisteringCustomers func(emailAddress, givenName, familyName string) (value.CustomerID, error)
