package customergrpc

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

	"github.com/AntonStoeckl/go-iddd/service/customer/application"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	register            application.ForRegisteringCustomers
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses
	changeEmailAddress  application.ForChangingCustomerEmailAddresses
	changeName          application.ForChangingCustomerNames
	delete              application.ForDeletingCustomers
	retrieveView        application.ForRetrievingCustomerViews
}

func NewCustomerServer(
	register application.ForRegisteringCustomers,
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses,
	changeEmailAddress application.ForChangingCustomerEmailAddresses,
	changeName application.ForChangingCustomerNames,
	delete application.ForDeletingCustomers,
	retrieveView application.ForRetrievingCustomerViews,
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
	ctx context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {

	command, err := commands.BuildRegisterCustomer(req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	if err := server.register(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &RegisterResponse{Id: command.CustomerID().ID()}, nil
}

func (server *customerServer) ConfirmEmailAddress(
	ctx context.Context,
	req *ConfirmEmailAddressRequest,
) (*empty.Empty, error) {

	command, err := commands.BuildConfirmCustomerEmailAddress(req.Id, req.ConfirmationHash)
	if err != nil {
		return nil, err
	}

	if err := server.confirmEmailAddress(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(
	ctx context.Context,
	req *ChangeEmailAddressRequest,
) (*empty.Empty, error) {

	command, err := commands.BuildChangeCustomerEmailAddress(req.Id, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	if err := server.changeEmailAddress(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeName(
	ctx context.Context,
	req *ChangeNameRequest,
) (*empty.Empty, error) {

	command, err := commands.BuildChangeCustomerName(req.Id, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	if err := server.changeName(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) Delete(
	ctx context.Context,
	req *DeleteRequest,
) (*empty.Empty, error) {

	command, err := commands.BuildDeleteCustomer(req.Id)
	if err != nil {
		return nil, err
	}

	if err := server.delete(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) RetrieveView(
	ctx context.Context,
	req *RetrieveViewRequest,
) (*RetrieveViewResponse, error) {

	id, err := values.BuildCustomerID(req.Id)
	if err != nil {
		return nil, err
	}

	view, err := server.retrieveView(id)
	if err != nil {
		return nil, err
	}

	response := &RetrieveViewResponse{
		EmailAddress:            view.EmailAddress,
		IsEmailAddressConfirmed: view.IsEmailAddressConfirmed,
		GivenName:               view.GivenName,
		FamilyName:              view.FamilyName,
		Version:                 uint64(view.Version),
	}

	_ = ctx // currently not used

	return response, nil
}
