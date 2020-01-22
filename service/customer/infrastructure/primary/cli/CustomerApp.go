package customercli

import (
	"fmt"
	"go-iddd/service/customer"
	"go-iddd/service/customer/application/domain/commands"

	"github.com/urfave/cli"
)

type CustomerApp struct {
	forRegisteringCustomers     customer.ForRegisteringCustomers
	forConfirmingEmailAddresses customer.ForConfirmingEmailAddresses
	forChangingEmailAddresses   customer.ForChangingEmailAddresses
}

func NewCustomerApp(
	forRegisteringCustomers customer.ForRegisteringCustomers,
	forConfirmingEmailAddresses customer.ForConfirmingEmailAddresses,
	forChangingEmailAddresses customer.ForChangingEmailAddresses,
) *CustomerApp {
	app := &CustomerApp{
		forRegisteringCustomers:     forRegisteringCustomers,
		forConfirmingEmailAddresses: forConfirmingEmailAddresses,
		forChangingEmailAddresses:   forChangingEmailAddresses,
	}

	return app
}

func (app *CustomerApp) Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "RegisterCustomer",
			Aliases:   []string{"rc"},
			Usage:     "Register a Customer",
			Action:    app.RegisterCustomer,
			ArgsUsage: "emailAddress givenName familyName",
		},
		{
			Name:      "ConfirmCustomerEmailAddress",
			Aliases:   []string{"cocea"},
			Usage:     "Confirm a Customer's emailAddress",
			Action:    app.ConfirmCustomerEmailAddress,
			ArgsUsage: "id emailAddress confirmationHash",
		},
		{
			Name:      "ChangeCustomerEmailAddress",
			Aliases:   []string{"chcea"},
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

	command, err := commands.NewRegister(emailAddress, givenName, familyName)
	if err != nil {
		return err
	}

	if err := app.forRegisteringCustomers.Register(command); err != nil {
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
	emailAddress := ctx.Args().Get(1)
	confirmationHash := ctx.Args().Get(2)

	command, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
	if err != nil {
		return err
	}

	if err := app.forConfirmingEmailAddresses.ConfirmEmailAddress(command); err != nil {
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

	command, err := commands.NewChangeEmailAddress(id, emailAddress)
	if err != nil {
		return err
	}

	if err := app.forChangingEmailAddresses.ChangeEmailAddress(command); err != nil {
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
