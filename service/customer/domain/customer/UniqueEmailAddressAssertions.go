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

func shouldAddUniqueEmailAddress(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
) UniqueEmailAddressAssertion {

	return UniqueEmailAddressAssertion{
		assertionType:     ShouldAddUniqueEmailAddress,
		customerID:        customerID,
		emailAddressToAdd: emailAddress,
	}
}

func shouldReplaceUniqueEmailAddress(
	newEmailAddress values.EmailAddress,
	currentEmailAddress values.EmailAddress,
) UniqueEmailAddressAssertion {

	return UniqueEmailAddressAssertion{
		assertionType:        ShouldReplaceUniqueEmailAddress,
		emailAddressToAdd:    newEmailAddress,
		emailAddressToRemove: currentEmailAddress,
	}
}

func shouldRemoveUniqueEmailAddress(
	emailAddress values.EmailAddress,
) UniqueEmailAddressAssertion {

	return UniqueEmailAddressAssertion{
		assertionType:        ShouldRemoveUniqueEmailAddress,
		emailAddressToRemove: emailAddress,
	}
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
				shouldAddUniqueEmailAddress(actualEvent.CustomerID(), actualEvent.EmailAddress()),
			)
		case events.CustomerEmailAddressChanged:
			specifications = append(
				specifications,
				shouldReplaceUniqueEmailAddress(actualEvent.EmailAddress(), actualEvent.PreviousEmailAddress()),
			)
		case events.CustomerDeleted:
			specifications = append(
				specifications,
				shouldRemoveUniqueEmailAddress(actualEvent.EmailAddress()),
			)
		}
	}

	return specifications
}
