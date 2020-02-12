package customergrpc

import (
	"context"
	"go-iddd/service/customer"
	"go-iddd/service/customer/application/domain/commands"

	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	register            customer.ForRegistering
	confirmEmailAddress customer.ForConfirmingEmailAddresses
	changeEmailAddress  customer.ForChangingEmailAddresses
}

func NewCustomerServer(
	register customer.ForRegistering,
	confirmEmailAddress customer.ForConfirmingEmailAddresses,
	changeEmailAddress customer.ForChangingEmailAddresses,
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

	command, err := commands.NewRegister(req.EmailAddress, req.GivenName, req.FamilyName)
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

	command, err := commands.NewConfirmEmailAddress(req.Id, req.EmailAddress, req.ConfirmationHash)
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

	command, err := commands.NewChangeEmailAddress(req.Id, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	if err := server.changeEmailAddress(command); err != nil {
		return nil, err
	}

	_ = ctx // currently not used

	return &empty.Empty{}, nil
}
