package customeraccounts

type ForConfirmingCustomerEmailAddresses func(customerID, confirmationHash string) error
