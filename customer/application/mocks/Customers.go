// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate by running `go generate` in customer/domain

package mocks

import domain "go-iddd/customer/domain"
import mock "github.com/stretchr/testify/mock"
import valueobjects "go-iddd/customer/domain/valueobjects"

// Customers is an autogenerated mock type for the Customers type
type Customers struct {
	mock.Mock
}

// FindBy provides a mock function with given fields: id
func (_m *Customers) FindBy(id *valueobjects.ID) (domain.Customer, error) {
	ret := _m.Called(id)

	var r0 domain.Customer
	if rf, ok := ret.Get(0).(func(*valueobjects.ID) domain.Customer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Customer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*valueobjects.ID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// New provides a mock function with given fields:
func (_m *Customers) New() domain.Customer {
	ret := _m.Called()

	var r0 domain.Customer
	if rf, ok := ret.Get(0).(func() domain.Customer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Customer)
		}
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *Customers) Save(_a0 domain.Customer) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Customer) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
