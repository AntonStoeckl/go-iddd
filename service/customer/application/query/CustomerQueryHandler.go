package query

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type CustomerQueryHandler struct {
	customerEvents ForReadingCustomerEventStreams
}

func NewCustomerQueryHandler(customerEvents ForReadingCustomerEventStreams) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		customerEvents: customerEvents,
	}
}

func (h *CustomerQueryHandler) CustomerViewByID(customerID string) (customer.View, error) {
	var err error
	var customerIDValue values.CustomerID
	wrapWithMsg := "customerQueryHandler.CustomerViewByID"

	if customerIDValue, err = values.BuildCustomerID(customerID); err != nil {
		return customer.View{}, errors.Wrap(err, wrapWithMsg)
	}

	eventStream, err := h.customerEvents.EventStreamFor(customerIDValue)
	if err != nil {
		return customer.View{}, errors.Wrap(err, wrapWithMsg)
	}

	customerView := customer.BuildViewFrom(eventStream)

	if customerView.IsDeleted {
		err := errors.New("customer not found")

		return customer.View{}, lib.MarkAndWrapError(err, lib.ErrNotFound, wrapWithMsg)
	}

	return customerView, nil
}
