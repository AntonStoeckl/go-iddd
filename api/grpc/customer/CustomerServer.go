package customer

import (
	"context"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
)

type customerServer struct{}

func NewCustomerServer() *customerServer {
	return &customerServer{}
}

func (*customerServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	id := values.GenerateID()

	_, err := commands.NewRegister(id.String(), req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{Id: id.String()}, nil
}
