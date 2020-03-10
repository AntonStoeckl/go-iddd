package readmodel

type CustomerQueryHandler struct {
	customerEvents ForReadingCustomerEvents
}

func NewCustomerQueryHandler(customerEvents ForReadingCustomerEvents) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		customerEvents: customerEvents,
	}
}
