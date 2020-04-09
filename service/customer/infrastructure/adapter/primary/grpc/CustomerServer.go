package customergrpc

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/service/customer/application"
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

	customerID, err := server.register(req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, MapToGRPCErrors(err)
	}

	_ = ctx // currently not used

	return &RegisterResponse{Id: customerID.String()}, nil
}

func (server *customerServer) ConfirmEmailAddress(
	ctx context.Context,
	req *ConfirmEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.confirmEmailAddress(req.Id, req.ConfirmationHash); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(
	ctx context.Context,
	req *ChangeEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.changeEmailAddress(req.Id, req.EmailAddress); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeName(
	ctx context.Context,
	req *ChangeNameRequest,
) (*empty.Empty, error) {

	if err := server.changeName(req.Id, req.GivenName, req.FamilyName); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) Delete(
	ctx context.Context,
	req *DeleteRequest,
) (*empty.Empty, error) {

	if err := server.delete(req.Id); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}

func (server *customerServer) RetrieveView(
	ctx context.Context,
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

	_ = ctx // currently not used

	return response, nil
}
