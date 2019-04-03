package command

import (
	"go-iddd/person/model/vo"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
}

func Test_RegisterTestSuite(t *testing.T) {
	tests := new(RegisterTestSuite)
	suite.Run(t, tests)
}

func (suite *RegisterTestSuite) Test_NewRegister() {
	// given
	id := vo.NewID("12345")
	emailAddress := vo.NewEmailAddress("foo@bar.com")
	name := vo.NewName("Anton", "Stöckl")

	// when
	command, err := NewRegister(id, emailAddress, name)

	// then
	suite.NoError(err)
	expectedCommandName := "Register"
	suite.Equal(expectedCommandName, command.CommandName(), "the CommandName should be %s", expectedCommandName)
	suite.Equal(id.ID(), command.Identifier(), "the Identifier should be %s", id.ID())
}

func (suite *RegisterTestSuite) Test_NewRegister_WithNilID() {
	// given
	var id vo.ID
	emailAddress := vo.NewEmailAddress("foo@bar.com")
	name := vo.NewName("Anton", "Stöckl")

	// when
	command, err := NewRegister(id, emailAddress, name)

	// then
	suite.Error(err, "it should fail because ID is nil")
	suite.Nil(command, "the command should be nil")
}

func (suite *RegisterTestSuite) Test_NewRegister_WithNilEmailAddress() {
	// given
	id := vo.NewID("12345")
	var emailAddress vo.EmailAddress
	name := vo.NewName("Anton", "Stöckl")

	// when
	command, err := NewRegister(id, emailAddress, name)

	// then
	suite.Error(err, "it should fail because EmailAddress is nil")
	suite.Nil(command, "the command should be nil")
}

func (suite *RegisterTestSuite) Test_NewRegister_WithNilName() {
	// given
	id := vo.NewID("12345")
	emailAddress := vo.NewEmailAddress("foo@bar.com")
	var name vo.Name

	// when
	command, err := NewRegister(id, emailAddress, name)

	// then
	suite.Error(err, "it should fail because Name is nil")
	suite.Nil(command, "the command should be nil")
}
