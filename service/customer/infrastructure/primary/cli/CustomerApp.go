package customercli

import (
	"fmt"
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"

	"github.com/urfave/cli"
)

type CustomerApp struct {
	register            application.ForRegisteringCustomers
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses
	changeEmailAddress  application.ForChangingCustomerEmailAddresses
}

func NewCustomerApp(
	register application.ForRegisteringCustomers,
	confirmEmailAddress application.ForConfirmingCustomerEmailAddresses,
	changeEmailAddress application.ForChangingCustomerEmailAddresses,
) *CustomerApp {
	app := &CustomerApp{
		register:            register,
		confirmEmailAddress: confirmEmailAddress,
		changeEmailAddress:  changeEmailAddress,
	}

	return app
}

func (app *CustomerApp) Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "Register",
			Aliases:   []string{"reg"},
			Usage:     "Register a Customer",
			Action:    app.RegisterCustomer,
			ArgsUsage: "emailAddress givenName familyName",
		},
		{
			Name:      "ConfirmCustomerEmailAddress",
			Aliases:   []string{"coea"},
			Usage:     "Confirm a Customer's emailAddress",
			Action:    app.ConfirmCustomerEmailAddress,
			ArgsUsage: "id confirmationHash",
		},
		{
			Name:      "ChangeCustomerEmailAddress",
			Aliases:   []string{"chea"},
			Usage:     "Change a Customer's emailAddress",
			Action:    app.ChangeCustomerEmailAddress,
			ArgsUsage: "id emailAddress",
		},
	}
}

func (app *CustomerApp) RegisterCustomer(ctx *cli.Context) error {
	emailAddress := ctx.Args().Get(0)
	givenName := ctx.Args().Get(1)
	familyName := ctx.Args().Get(2)

	command, err := commands.BuildRegisterCustomer(emailAddress, givenName, familyName)
	if err != nil {
		return err
	}

	if err := app.register(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"Customer registered with id '%s'\n",
		command.CustomerID().ID(),
	)

	return nil
}

func (app *CustomerApp) ConfirmCustomerEmailAddress(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	confirmationHash := ctx.Args().Get(1)

	command, err := commands.BuildConfirmCustomerEmailAddress(id, confirmationHash)
	if err != nil {
		return err
	}

	if err := app.confirmEmailAddress(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"successfully confirmed the emailAddress of Customer with id '%s'\n",
		command.CustomerID().ID(),
	)

	return nil
}

func (app *CustomerApp) ChangeCustomerEmailAddress(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	emailAddress := ctx.Args().Get(1)

	command, err := commands.BuildChangeCustomerEmailAddress(id, emailAddress)
	if err != nil {
		return err
	}

	if err := app.changeEmailAddress(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"successfully changed the emailAddress of Customer with id '%s' to '%s\n",
		command.CustomerID().ID(),
		command.EmailAddress().EmailAddress(),
	)

	return nil
}
