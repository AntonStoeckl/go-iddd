package customer

import "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

func IsMatchingConfirmationHash(current values.ConfirmationHash, supplied values.ConfirmationHash) bool {
	return current.Equals(supplied)
}
