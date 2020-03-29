package application

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

type ForRegisteringCustomers func(emailAddress, givenName, familyName string) (values.CustomerID, error)
