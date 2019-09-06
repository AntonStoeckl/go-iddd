package customer

import (
	"context"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type customerServer struct {
	commandHandler shared.CommandHandler
}

func NewCustomerServer(commandHandeler shared.CommandHandler) *customerServer {
	return &customerServer{commandHandler: commandHandeler}
}

func (customerServer *customerServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	id := values.GenerateID()

	register, err := commands.NewRegister(id.String(), req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	if err := customerServer.commandHandler.Handle(register); err != nil {
		return nil, err
	}

	return &RegisterResponse{Id: id.String()}, nil
}
