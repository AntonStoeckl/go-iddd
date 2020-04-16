package customergrpc

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts"
	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	register            customeraccounts.ForRegisteringCustomers
	confirmEmailAddress customeraccounts.ForConfirmingCustomerEmailAddresses
	changeEmailAddress  customeraccounts.ForChangingCustomerEmailAddresses
	changeName          customeraccounts.ForChangingCustomerNames
	delete              customeraccounts.ForDeletingCustomers
	retrieveView        customeraccounts.ForRetrievingCustomerViews
}

func NewCustomerServer(
	register customeraccounts.ForRegisteringCustomers,
	confirmEmailAddress customeraccounts.ForConfirmingCustomerEmailAddresses,
	changeEmailAddress customeraccounts.ForChangingCustomerEmailAddresses,
	changeName customeraccounts.ForChangingCustomerNames,
	delete customeraccounts.ForDeletingCustomers,
	retrieveView customeraccounts.ForRetrievingCustomerViews,
) *customerServer {
	server := &customerServer{
		register:            register,
		confirmEmailAddress: confirmEmailAddress,
		changeEmailAddress:  changeEmailAddress,
		changeName:          changeName,
		delete:              delete,
		retrieveView:        retrieveView,
	}

	return server
}

func (server *customerServer) Register(
	_ context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {

	customerID, err := server.register(req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &RegisterResponse{Id: customerID.String()}, nil
}

func (server *customerServer) ConfirmEmailAddress(
	_ context.Context,
	req *ConfirmEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.confirmEmailAddress(req.Id, req.ConfirmationHash); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(
	_ context.Context,
	req *ChangeEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.changeEmailAddress(req.Id, req.EmailAddress); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeName(
	_ context.Context,
	req *ChangeNameRequest,
) (*empty.Empty, error) {

	if err := server.changeName(req.Id, req.GivenName, req.FamilyName); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) Delete(
	_ context.Context,
	req *DeleteRequest,
) (*empty.Empty, error) {

	if err := server.delete(req.Id); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) RetrieveView(
	_ context.Context,
	req *RetrieveViewRequest,
) (*RetrieveViewResponse, error) {

	view, err := server.retrieveView(req.Id)
	if err != nil {
		return nil, MapToGRPCErrors(err)
	}

	response := &RetrieveViewResponse{
		EmailAddress:            view.EmailAddress,
		IsEmailAddressConfirmed: view.IsEmailAddressConfirmed,
		GivenName:               view.GivenName,
		FamilyName:              view.FamilyName,
		Version:                 uint64(view.Version),
	}

	return response, nil
}
