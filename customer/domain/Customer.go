package domain

import (
	"errors"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

//go:generate mockery -name Customer -output ../application/mocks -outpkg mocks -note "Regenerate by running `go generate` in customer/domain"

type Customer interface {
	Apply(cmd shared.Command) error
}

type customer struct {
	id           valueobjects.ID
	emailAddress valueobjects.ConfirmableEmailAddress
	personName   valueobjects.PersonName
	isRegistered bool
}

func NewUnregisteredCustomer() *customer {
	return &customer{}
}

func (customer *customer) Apply(command shared.Command) error {
	var err error

	if err := customer.assertCustomerIsInValidState(command); err != nil {
		return err
	}

	/*** All methods to apply the commands to the Customer are located in the Commands itself ***/

	switch command := command.(type) {
	case Register:
		customer.register(command)
	case ConfirmEmailAddress:
		err = customer.confirmEmailAddress(command)
	case nil:
		err = errors.New("customer - nil command applied")
	default:
		err = errors.New("customer - unknown command applied")
	}

	return err
}

func (customer *customer) assertCustomerIsInValidState(command shared.Command) error {
	switch command.(type) {
	case Register:
		if customer.isRegistered {
			return errors.New("customer - was already registered")
		}
	default:
		if !customer.isRegistered {
			return errors.New("customer - was not registered yet")
		}

		if customer.id == nil {
			return errors.New("customer - was registered but has no id")
		}

		if customer.emailAddress == nil {
			return errors.New("customer - was registered but has no emailAddress")
		}

		if customer.personName == nil {
			return errors.New("customer - was registered but has no personName")
		}
	}

	return nil
}
