package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

const (
	ShouldAddUniqueEmailAddress = iota
	ShouldReplaceUniqueEmailAddress
	ShouldRemoveUniqueEmailAddress
)

type ForBuildingUniqueEmailAddressAssertions func(recordedEvents ...es.DomainEvent) UniqueEmailAddressAssertions

type UniqueEmailAddressAssertion struct {
	desiredAction     int
	customerID        value.CustomerID
	emailAddressToAdd value.UnconfirmedEmailAddress
}

type UniqueEmailAddressAssertions []UniqueEmailAddressAssertion

func (spec UniqueEmailAddressAssertion) DesiredAction() int {
	return spec.desiredAction
}

func (spec UniqueEmailAddressAssertion) CustomerID() value.CustomerID {
	return spec.customerID
}

func (spec UniqueEmailAddressAssertion) EmailAddressToAdd() value.UnconfirmedEmailAddress {
	return spec.emailAddressToAdd
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
					desiredAction:     ShouldReplaceUniqueEmailAddress,
					customerID:        actualEvent.CustomerID(),
					emailAddressToAdd: actualEvent.EmailAddress(),
				},
			)
		case domain.CustomerDeleted:
			specifications = append(
				specifications,
				UniqueEmailAddressAssertion{
					desiredAction: ShouldRemoveUniqueEmailAddress,
					customerID:    actualEvent.CustomerID(),
				},
			)
		}
	}

	return specifications
}
