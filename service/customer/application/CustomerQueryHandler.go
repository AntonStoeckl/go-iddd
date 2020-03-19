package application

import (
	"go-iddd/service/customer/application/domain/customer"
	"go-iddd/service/customer/application/domain/values"

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

func (h *CustomerQueryHandler) CustomerViewByID(customerID values.CustomerID) (customer.View, error) {
	eventStream, err := h.customerEvents.EventStreamFor(customerID)
	if err != nil {
		return customer.View{}, errors.Wrap(err, "customerQueryHandler.CustomerViewByID")
	}

	customerView := customer.BuildViewFrom(eventStream)

	return customerView, nil
}
