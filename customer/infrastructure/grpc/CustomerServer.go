package customergrpc

import (
	"context"
	"go-iddd/customer"
	"go-iddd/customer/application"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"

	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	forRegisteringCustomers     customer.ForRegisteringCustomers
	forConfirmingEmailAddresses customer.ForConfirmingEmailAddresses
	forChangingEmailAddresses   customer.ForChangingEmailAddresses
}

func NewCustomerServer(commandHandler *application.CommandHandler) *customerServer {
	server := &customerServer{
		forRegisteringCustomers:     commandHandler,
		forConfirmingEmailAddresses: commandHandler,
		forChangingEmailAddresses:   commandHandler,
	}

	return server
}

func (server *customerServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	id := values.GenerateCustomerID()

	command, err := commands.NewRegister(id.ID(), req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	if err := server.forRegisteringCustomers.Register(command); err != nil {
		return nil, err
	}

	return &RegisterResponse{Id: id.ID()}, nil
}

func (server *customerServer) ConfirmEmailAddress(ctx context.Context, req *ConfirmEmailAddressRequest) (*empty.Empty, error) {
	command, err := commands.NewConfirmEmailAddress(req.Id, req.EmailAddress, req.ConfirmationHash)
	if err != nil {
		return nil, err
	}

	if err := server.forConfirmingEmailAddresses.ConfirmEmailAddress(command); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(ctx context.Context, req *ChangeEmailAddressRequest) (*empty.Empty, error) {
	command, err := commands.NewChangeEmailAddress(req.Id, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	if err := server.forChangingEmailAddresses.ChangeEmailAddress(command); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
