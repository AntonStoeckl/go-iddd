package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

type CustomerQueryHandler struct {
	retrieveCustomerEventStream ForRetrievingCustomerEventStreams
}

func NewCustomerQueryHandler(retrieveCustomerEventStream ForRetrievingCustomerEventStreams) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		retrieveCustomerEventStream: retrieveCustomerEventStream,
	}
}

func (h *CustomerQueryHandler) CustomerViewByID(customerID string) (customer.View, error) {
	var err error
	var customerIDValue value.CustomerID
	wrapWithMsg := "customerQueryHandler.CustomerViewByID"

	if customerIDValue, err = value.BuildCustomerID(customerID); err != nil {
		return customer.View{}, errors.Wrap(err, wrapWithMsg)
	}

	eventStream, err := h.retrieveCustomerEventStream(customerIDValue)
	if err != nil {
		return customer.View{}, errors.Wrap(err, wrapWithMsg)
	}

	customerView := customer.BuildViewFrom(eventStream)

	if customerView.IsDeleted {
		err := errors.New("customer not found")

		return customer.View{}, shared.MarkAndWrapError(err, shared.ErrNotFound, wrapWithMsg)
	}

	return customerView, nil
}
