package customeraccounts

type ForChangingCustomerNames func(customerID, givenName, familyName string) error
