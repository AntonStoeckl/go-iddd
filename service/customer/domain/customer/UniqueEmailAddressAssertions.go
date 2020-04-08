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

type UniqueEmailAddressAssertion struct {
	assertionType        int
	customerID           values.CustomerID
	emailAddressToAdd    values.EmailAddress
	emailAddressToRemove values.EmailAddress
}

func (spec UniqueEmailAddressAssertion) AssertionType() int {
	return spec.assertionType
}

func (spec UniqueEmailAddressAssertion) CustomerID() values.CustomerID {
	return spec.customerID
}

func (spec UniqueEmailAddressAssertion) EmailAddressToAdd() values.EmailAddress {
	return spec.emailAddressToAdd
}

func (spec UniqueEmailAddressAssertion) EmailAddressToRemove() values.EmailAddress {
	return spec.emailAddressToRemove
}

func BuildUniqueEmailAddressAssertionsFrom(recordedEvents es.DomainEvents) []UniqueEmailAddressAssertion {
	var specifications []UniqueEmailAddressAssertion

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					assertionType:     ShouldAddUniqueEmailAddress,
					customerID:        actualEvent.CustomerID(),
					emailAddressToAdd: actualEvent.EmailAddress(),
				},
			)
		case events.CustomerEmailAddressChanged:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					assertionType:        ShouldReplaceUniqueEmailAddress,
					emailAddressToAdd:    actualEvent.EmailAddress(),
					emailAddressToRemove: actualEvent.PreviousEmailAddress(),
				},
			)
		case events.CustomerDeleted:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					assertionType:        ShouldRemoveUniqueEmailAddress,
					emailAddressToRemove: actualEvent.EmailAddress(),
				},
			)
		}
	}

	return specifications
}
