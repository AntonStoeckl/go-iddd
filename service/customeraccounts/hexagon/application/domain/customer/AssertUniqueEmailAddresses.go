package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

const (
	ShouldAddUniqueEmailAddress = iota
	ShouldReplaceUniqueEmailAddress
	ShouldRemoveUniqueEmailAddress
)

type ForBuildingUniqueEmailAddressAssertions func(recordedEvents ...es.DomainEvent) UniqueEmailAddressAssertions

type UniqueEmailAddressAssertion struct {
	desiredAction        int
	customerID           value.CustomerID
	emailAddressToAdd    value.EmailAddress
	emailAddressToRemove value.EmailAddress
}

type UniqueEmailAddressAssertions []UniqueEmailAddressAssertion

func (spec UniqueEmailAddressAssertion) DesiredAction() int {
	return spec.desiredAction
}

func (spec UniqueEmailAddressAssertion) CustomerID() value.CustomerID {
	return spec.customerID
}

func (spec UniqueEmailAddressAssertion) EmailAddressToAdd() value.EmailAddress {
	return spec.emailAddressToAdd
}

func (spec UniqueEmailAddressAssertion) EmailAddressToRemove() value.EmailAddress {
	return spec.emailAddressToRemove
}

func BuildUniqueEmailAddressAssertions(recordedEvents ...es.DomainEvent) UniqueEmailAddressAssertions {
	var specifications UniqueEmailAddressAssertions

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					desiredAction:     ShouldAddUniqueEmailAddress,
					customerID:        actualEvent.CustomerID(),
					emailAddressToAdd: actualEvent.EmailAddress(),
				},
			)
		case domain.CustomerEmailAddressChanged:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					desiredAction:        ShouldReplaceUniqueEmailAddress,
					emailAddressToAdd:    actualEvent.EmailAddress(),
					emailAddressToRemove: actualEvent.PreviousEmailAddress(),
				},
			)
		case domain.CustomerDeleted:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					desiredAction:        ShouldRemoveUniqueEmailAddress,
					emailAddressToRemove: actualEvent.EmailAddress(),
				},
			)
		}
	}

	return specifications
}
