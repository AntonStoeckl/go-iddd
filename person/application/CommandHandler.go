package application

import (
	"go-iddd/person/model"
	"go-iddd/person/model/command"
	"go-iddd/shared"
)

type commandHandler struct {
	persons model.Persons
}

func NewCommandHandler(persons model.Persons) *commandHandler {
	return &commandHandler{persons: persons}
}

func (handler *commandHandler) Handle(commandx shared.Command) error {
	var err error

	switch commandx := commandx.(type) {
	case command.Register:
		err = handler.Register(commandx)
	}

	return err
}

func (handler *commandHandler) Register(command command.Register) error {
	person := model.NewPerson()
	person.Register(command)
	err := handler.persons.Save(person)

	return err
}
