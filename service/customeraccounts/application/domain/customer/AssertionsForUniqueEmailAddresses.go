package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

const (
	ShouldAddUniqueEmailAddress = iota
	ShouldReplaceUniqueEmailAddress
	ShouldRemoveUniqueEmailAddress
)

type AssertionForUniqueEmailAddresses struct {
	desiredAction        int
	customerID           value.CustomerID
	emailAddressToAdd    value.EmailAddress
	emailAddressToRemove value.EmailAddress
}

type AssertionsForUniqueEmailAddresses []AssertionForUniqueEmailAddresses

func (spec AssertionForUniqueEmailAddresses) DesiredAction() int {
	return spec.desiredAction
}

func (spec AssertionForUniqueEmailAddresses) CustomerID() value.CustomerID {
	return spec.customerID
}

func (spec AssertionForUniqueEmailAddresses) EmailAddressToAdd() value.EmailAddress {
	return spec.emailAddressToAdd
}

func (spec AssertionForUniqueEmailAddresses) EmailAddressToRemove() value.EmailAddress {
	return spec.emailAddressToRemove
}

func BuildAssertionsForUniqueEmailAddresses(recordedEvents es.RecordedEvents) AssertionsForUniqueEmailAddresses {
	var specifications AssertionsForUniqueEmailAddresses

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			specifications = append(
				specifications,
				AssertionForUniqueEmailAddresses{
					desiredAction:     ShouldAddUniqueEmailAddress,
					customerID:        actualEvent.CustomerID(),
					emailAddressToAdd: actualEvent.EmailAddress(),
				},
			)
		case domain.CustomerEmailAddressChanged:
			specifications = append(
				specifications,
				AssertionForUniqueEmailAddresses{
					desiredAction:        ShouldReplaceUniqueEmailAddress,
					emailAddressToAdd:    actualEvent.EmailAddress(),
					emailAddressToRemove: actualEvent.PreviousEmailAddress(),
				},
			)
		case domain.CustomerDeleted:
			specifications = append(
				specifications,
				AssertionForUniqueEmailAddresses{
					desiredAction:        ShouldRemoveUniqueEmailAddress,
					emailAddressToRemove: actualEvent.EmailAddress(),
				},
			)
		}
	}

	return specifications
}
