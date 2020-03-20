package customergrpc

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/service/customer/application"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	register            application.ForRegisteringCustomers
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses
	changeEmailAddress  application.ForChangingCustomerEmailAddresses
}

func NewCustomerServer(
	register application.ForRegisteringCustomers,
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses,
	changeEmailAddress application.ForChangingCustomerEmailAddresses,
) *customerServer {
	server := &customerServer{
		register:            register,
		confirmEmailAddress: confirmEmailAddress,
		changeEmailAddress:  changeEmailAddress,
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
