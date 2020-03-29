package application

type ForConfirmingCustomerEmailAddresses func(customerID, confirmationHash string) error
