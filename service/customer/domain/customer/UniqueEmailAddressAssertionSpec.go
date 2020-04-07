package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

const (
	AddUniqueEmailAddress = iota
	ReplaceUniqueEmailAddress
	RemoveUniqueEmailAddress
)

type UniqueEmailAddressAssertionSpec struct {
	trackingOperationType int
	customerID            values.CustomerID
	emailAddressToAdd     values.EmailAddress
	emailAddressToRemove  values.EmailAddress
}

func shouldAddUniqueEmailAddress(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
) UniqueEmailAddressAssertionSpec {

	return UniqueEmailAddressAssertionSpec{
		trackingOperationType: AddUniqueEmailAddress,
		customerID:            customerID,
		emailAddressToAdd:     emailAddress,
	}
}

func shouldReplaceUniqueEmailAddress(
	newEmailAddress values.EmailAddress,
	currentEmailAddress values.EmailAddress,
) UniqueEmailAddressAssertionSpec {

	return UniqueEmailAddressAssertionSpec{
		trackingOperationType: ReplaceUniqueEmailAddress,
		emailAddressToAdd:     newEmailAddress,
		emailAddressToRemove:  currentEmailAddress,
	}
}

func shouldRemoveUniqueEmailAddress(
	emailAddress values.EmailAddress,
) UniqueEmailAddressAssertionSpec {

	return UniqueEmailAddressAssertionSpec{
		trackingOperationType: RemoveUniqueEmailAddress,
		emailAddressToRemove:  emailAddress,
	}
}

func (spec UniqueEmailAddressAssertionSpec) TrackingOperationType() int {
	return spec.trackingOperationType
}

func (spec UniqueEmailAddressAssertionSpec) CustomerID() values.CustomerID {
	return spec.customerID
}

func (spec UniqueEmailAddressAssertionSpec) EmailAddressToAdd() values.EmailAddress {
	return spec.emailAddressToAdd
}

func (spec UniqueEmailAddressAssertionSpec) EmailAddressToRemove() values.EmailAddress {
	return spec.emailAddressToRemove
}

func DeriveUniqueEmailAddressAssertionSpecsFrom(recordedEvents es.DomainEvents) []UniqueEmailAddressAssertionSpec {
	var specifications []UniqueEmailAddressAssertionSpec

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
