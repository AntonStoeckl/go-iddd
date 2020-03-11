package readmodel

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/readmodel/domain/customer/queries"
	"go-iddd/service/customer/application/writemodel/domain/customer/values"

	"github.com/cockroachdb/errors"
)

type CustomerQueryHandler struct {
	customerEvents ForReadingCustomerEvents
}

func NewCustomerQueryHandler(customerEvents ForReadingCustomerEvents) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		customerEvents: customerEvents,
	}
}

func (h *CustomerQueryHandler) CustomerViewByID(query queries.CustomerByID) (customer.View, error) {
	eventStream, err := h.customerEvents.EventStreamFor(values.RebuildCustomerID(query.CustomerID))
	if err != nil {
		return customer.View{}, errors.Wrap(err, "customerQueryHandler.CustomerViewByID")
	}

	customerView := customer.BuildViewFrom(eventStream)

	return customerView, nil
}
