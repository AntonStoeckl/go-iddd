package customer

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

func HasSuppliedMatchingConfirmationHash(current values.ConfirmationHash, supplied values.ConfirmationHash) bool {
	return current.Equals(supplied)
}
