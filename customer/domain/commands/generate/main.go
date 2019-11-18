// +build generator

//go:generate go run main.go

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/joncalhoun/pipe"
)

const mode = "format" // "format", "noformat", "stdout"

type Input struct {
	FieldName string
	DataType  string
	Valid     string
	Invalid   string
}

type Field struct {
	FieldName    string
	DataType     string
	ValueFactory string
	Input        []Input
}

type Command struct {
	CommandType string
	Fields      []Field
}

var commands = []Command{
	{
		CommandType: "Register",
		Fields: []Field{
			{
				FieldName:    "customerID",
				DataType:     "*values.CustomerID",
				ValueFactory: "values.CustomerIDFrom(customerID)",
				Input: []Input{
					{FieldName: "customerID", DataType: "string", Valid: `"64bcf656-da30-4f5a-b0b5-aead60965aa3"`, Invalid: `""`},
				},
			},
			{
				FieldName:    "emailAddress",
				DataType:     "*values.EmailAddress",
				ValueFactory: "values.EmailAddressFrom(emailAddress)",
				Input: []Input{
					{FieldName: "emailAddress", DataType: "string", Valid: `"john@doe.com"`, Invalid: `""`},
				},
			},
			{
				FieldName:    "personName",
				DataType:     "*values.PersonName",
				ValueFactory: "values.PersonNameFrom(givenName, familyName)",
				Input: []Input{
					{FieldName: "givenName", DataType: "string", Valid: `"John"`, Invalid: `""`},
					{FieldName: "familyName", DataType: "string", Valid: `"Doe"`, Invalid: `""`},
				},
			},
		},
	},
	{
		CommandType: "ConfirmEmailAddress",
		Fields: []Field{
			{
				FieldName:    "customerID",
				DataType:     "*values.CustomerID",
				ValueFactory: "values.CustomerIDFrom(customerID)",
				Input: []Input{
					{FieldName: "customerID", DataType: "string", Valid: `"64bcf656-da30-4f5a-b0b5-aead60965aa3"`, Invalid: `""`},
				},
			},
			{
				FieldName:    "emailAddress",
				DataType:     "*values.EmailAddress",
				ValueFactory: "values.EmailAddressFrom(emailAddress)",
				Input: []Input{
					{FieldName: "emailAddress", DataType: "string", Valid: `"john@doe.com"`, Invalid: `""`},
				},
			},
			{
				FieldName:    "confirmationHash",
				DataType:     "*values.ConfirmationHash",
				ValueFactory: "values.ConfirmationHashFrom(confirmationHash)",
				Input: []Input{
					{FieldName: "confirmationHash", DataType: "string", Valid: `"secret_hash"`, Invalid: `""`},
				},
			},
		},
	},
	{
		CommandType: "ChangeEmailAddress",
		Fields: []Field{
			{
				FieldName:    "customerID",
				DataType:     "*values.CustomerID",
				ValueFactory: "values.CustomerIDFrom(customerID)",
				Input: []Input{
					{FieldName: "customerID", DataType: "string", Valid: `"64bcf656-da30-4f5a-b0b5-aead60965aa3"`, Invalid: `""`},
				},
			},
			{
				FieldName:    "emailAddress",
				DataType:     "*values.EmailAddress",
				ValueFactory: "values.EmailAddressFrom(emailAddress)",
				Input: []Input{
					{FieldName: "emailAddress", DataType: "string", Valid: `"john@doe.com"`, Invalid: `""`},
				},
			},
		},
	},
}

type Config struct {
	RelativeOutputPath string
	Commands           []Command
}

var config = Config{
	RelativeOutputPath: "..",
}

func main() {
	generateCommands()
	generateTestsForCommands()
}

func generateCommands() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	lcFirst := func(s string) string {
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToLower(r)) + s[n:]
	}

	for _, command := range commands {
		t := template.New(command.CommandType)

		t = t.Funcs(
			template.FuncMap{
				"methodName":  methodName,
				"commandName": t.Name,
				"lcFirst":     lcFirst,
			},
		)

		t, err = t.Parse(commandTemplate)
		die(err)

		switch mode {
		case "format":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(command.CommandType) + ".go")
			die(err)

			rc, wc, _ := pipe.Commands(
				exec.Command("gofmt"),
				exec.Command("goimports"),
			)

			err = t.Execute(wc, command)
			die(err)

			err = wc.Close()
			die(err)

			_, err = io.Copy(outFile, rc)
			die(err)
		case "noformat":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(command.CommandType) + ".go")
			die(err)

			err = t.Execute(outFile, command)
			die(err)
		case "stdout":
			err = t.Execute(os.Stdout, command)
			die(err)
		}
	}
}

func generateTestsForCommands() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	lcFirst := func(s string) string {
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToLower(r)) + s[n:]
	}

	for _, event := range commands {
		t := template.New(event.CommandType)

		t = t.Funcs(
			template.FuncMap{
				"methodName":  methodName,
				"commandName": t.Name,
				"lcFirst":     lcFirst,
			},
		)

		t, err = t.Parse(testTemplate)
		die(err)

		switch mode {
		case "format":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.CommandType) + "_test.go")
			die(err)

			rc, wc, _ := pipe.Commands(
				exec.Command("gofmt"),
				exec.Command("goimports"),
			)

			err = t.Execute(wc, event)
			die(err)

			err = wc.Close()
			die(err)

			_, err = io.Copy(outFile, rc)
			die(err)
		case "noformat":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.CommandType) + "_test.go")
			die(err)

			err = t.Execute(outFile, event)
			die(err)
		case "stdout":
			err = t.Execute(os.Stdout, event)
			die(err)
		}
	}
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var commandTemplate = `
{{$commandVar := lcFirst commandName}}
// Code generated by generate/main.go. DO NOT EDIT.

package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
)

type {{commandName}} struct {
	{{range .Fields}}{{.FieldName}} {{.DataType}}
	{{end -}}
}

/*** Factory Method ***/

func New{{commandName}}(
	{{range .Fields}}
	{{- range .Input}}{{.FieldName}} {{.DataType}},
	{{end -}}
	{{end -}}
) (*{{commandName}}, error) {

	{{range .Fields}}
	{{.FieldName}}Value, err := {{.ValueFactory}}
	if err != nil {
		return nil, err
	}
	{{end}}

	{{$commandVar}} := &{{commandName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}}Value,
		{{end -}}
	}

	return {{$commandVar}}, nil
}

/*** Getter Methods ***/

{{range .Fields}}
func ({{$commandVar}} *{{commandName}}) {{methodName .DataType}}() {{.DataType}} {
	return {{$commandVar}}.{{.FieldName}}
}
{{end}}

/*** Implement shared.Command ***/

func ({{$commandVar}} *{{commandName}}) AggregateID() shared.IdentifiesAggregates {
	return {{$commandVar}}.customerID
}

func ({{$commandVar}} *{{commandName}}) CommandName() string {
	commandType := reflect.TypeOf({{$commandVar}}).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
`

var testTemplate = `
{{$commandVar := lcFirst commandName}}
{{$fields := .Fields}}
// Code generated by generate/main.go. DO NOT EDIT.

package commands_test

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNew{{commandName}}(t *testing.T) {
	Convey("Given valid input", t, func() {
		{{range .Fields}}{{range .Input}}{{.FieldName}} := {{.Valid}}
		{{end}}{{end}}

		Convey("When a new {{commandName}} command is created", func() {
			{{$commandVar}}, err := commands.New{{commandName}}({{range .Fields}}{{range .Input}}{{.FieldName}}, {{end}}{{end}})

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So({{$commandVar}}, ShouldHaveSameTypeAs, (*commands.{{commandName}})(nil))
			})
		})

		{{range .Fields}}
		{{range .Input}}
		Convey("Given that {{.FieldName}} is invalid", func() {
			{{.FieldName}} = {{.Invalid}}
			conveyNew{{commandName}}WithInvalidInput({{range $fields}}{{range .Input}}{{.FieldName}}, {{end}}{{end}})
		})
		{{end -}}
		{{end -}}
	})
}

func conveyNew{{commandName}}WithInvalidInput(
	{{range .Fields}}
	{{- range .Input}}{{.FieldName}} {{.DataType}},
	{{end -}}
	{{end -}}
) {

	Convey("When a new {{commandName}} command is created", func() {
		{{$commandVar}}, err := commands.New{{commandName}}({{range .Fields}}{{range .Input}}{{.FieldName}}, {{end}}{{end}})

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
			So({{$commandVar}}, ShouldBeNil)
		})
	})
}

func Test{{commandName}}ExposesExpectedValues(t *testing.T) {
	Convey("Given a {{commandName}} command", t, func() {
		{{range .Fields}}{{range .Input}}{{.FieldName}} := {{.Valid}}
		{{end}}{{end}}

		{{range .Fields}}{{.FieldName}}Value, err := {{.ValueFactory}}
		So(err, ShouldBeNil)
		{{end}}

		{{$commandVar}}, err := commands.New{{commandName}}({{range .Fields}}{{range .Input}}{{.FieldName}}, {{end}}{{end}})
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			{{range .Fields}}So({{.FieldName}}Value.Equals({{$commandVar}}.{{methodName .DataType}}()), ShouldBeTrue)
			{{end -}}
			So({{$commandVar}}.CommandName(), ShouldEqual, "{{commandName}}")
			So(customerIDValue.Equals({{$commandVar}}.AggregateID()), ShouldBeTrue)
		})
	})
}
`
