package customergrpc

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	customergrpcproto "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	register            hexagon.ForRegisteringCustomers
	confirmEmailAddress hexagon.ForConfirmingCustomerEmailAddresses
	changeEmailAddress  hexagon.ForChangingCustomerEmailAddresses
	changeName          hexagon.ForChangingCustomerNames
	delete              hexagon.ForDeletingCustomers
	retrieveView        hexagon.ForRetrievingCustomerViews
}

func NewCustomerServer(
	register hexagon.ForRegisteringCustomers,
	confirmEmailAddress hexagon.ForConfirmingCustomerEmailAddresses,
	changeEmailAddress hexagon.ForChangingCustomerEmailAddresses,
	changeName hexagon.ForChangingCustomerNames,
	delete hexagon.ForDeletingCustomers, //nolint:gocritic // false positive (shadowing of predeclared identifier: delete)
	retrieveView hexagon.ForRetrievingCustomerViews,
) customergrpcproto.CustomerServer {
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
	req *customergrpcproto.RegisterRequest,
) (*customergrpcproto.RegisterResponse, error) {

	customerIDValue := value.GenerateCustomerID()

	if err := server.register(customerIDValue, req.EmailAddress, req.GivenName, req.FamilyName); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &customergrpcproto.RegisterResponse{Id: customerIDValue.String()}, nil
}

func (server *customerServer) ConfirmEmailAddress(
	_ context.Context,
	req *customergrpcproto.ConfirmEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.confirmEmailAddress(req.Id, req.ConfirmationHash); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(
	_ context.Context,
	req *customergrpcproto.ChangeEmailAddressRequest,
) (*empty.Empty, error) {

	if err := server.changeEmailAddress(req.Id, req.EmailAddress); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeName(
	_ context.Context,
	req *customergrpcproto.ChangeNameRequest,
) (*empty.Empty, error) {

	if err := server.changeName(req.Id, req.GivenName, req.FamilyName); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) Delete(
	_ context.Context,
	req *customergrpcproto.DeleteRequest,
) (*empty.Empty, error) {

	if err := server.delete(req.Id); err != nil {
		return nil, MapToGRPCErrors(err)
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) RetrieveView(
	_ context.Context,
	req *customergrpcproto.RetrieveViewRequest,
) (*customergrpcproto.RetrieveViewResponse, error) {

	view, err := server.retrieveView(req.Id)
	if err != nil {
		return nil, MapToGRPCErrors(err)
	}

	response := &customergrpcproto.RetrieveViewResponse{
		EmailAddress:            view.EmailAddress,
		IsEmailAddressConfirmed: view.IsEmailAddressConfirmed,
		GivenName:               view.GivenName,
		FamilyName:              view.FamilyName,
		Version:                 uint64(view.Version),
	}

	return response, nil
}
