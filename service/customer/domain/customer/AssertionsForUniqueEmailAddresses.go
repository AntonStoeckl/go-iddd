package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

const (
	ShouldAddUniqueEmailAddress = iota
	ShouldReplaceUniqueEmailAddress
	ShouldRemoveUniqueEmailAddress
)

type AssertionForUniqueEmailAddresses struct {
	desiredAction        int
	customerID           values.CustomerID
	emailAddressToAdd    values.EmailAddress
	emailAddressToRemove values.EmailAddress
}

type AssertionsForUniqueEmailAddresses []AssertionForUniqueEmailAddresses

func (spec AssertionForUniqueEmailAddresses) DesiredAction() int {
	return spec.desiredAction
}

func (spec AssertionForUniqueEmailAddresses) CustomerID() values.CustomerID {
	return spec.customerID
}

func (spec AssertionForUniqueEmailAddresses) EmailAddressToAdd() values.EmailAddress {
	return spec.emailAddressToAdd
}

func (spec AssertionForUniqueEmailAddresses) EmailAddressToRemove() values.EmailAddress {
	return spec.emailAddressToRemove
}

func BuildAssertionsForUniqueEmailAddresses(recordedEvents es.RecordedEvents) AssertionsForUniqueEmailAddresses {
	var specifications AssertionsForUniqueEmailAddresses

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			specifications = append(
				specifications,
				AssertionForUniqueEmailAddresses{
					desiredAction:     ShouldAddUniqueEmailAddress,
					customerID:        actualEvent.CustomerID(),
					emailAddressToAdd: actualEvent.EmailAddress(),
				},
			)
		case events.CustomerEmailAddressChanged:
			specifications = append(
				specifications,
				AssertionForUniqueEmailAddresses{
					desiredAction:        ShouldReplaceUniqueEmailAddress,
					emailAddressToAdd:    actualEvent.EmailAddress(),
					emailAddressToRemove: actualEvent.PreviousEmailAddress(),
				},
			)
		case events.CustomerDeleted:
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
