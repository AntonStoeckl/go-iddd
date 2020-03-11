package queries

type CustomerByID struct {
	CustomerID string
}

func BuildCustomerByID(customerID string) CustomerByID {
	return CustomerByID{CustomerID: customerID}
}
