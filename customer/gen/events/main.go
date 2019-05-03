//go:generate go run main.go

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type Field struct {
	FieldName string
	DataType  string
}

type Event struct {
	EventType        string
	EventFactory     string
	AggregateFactory string
	Fields           []Field
}

var events = []Event{
	{
		EventType:    "registered",
		EventFactory: "ItWasRegistered",
		Fields: []Field{
			{FieldName: "id", DataType: "valueobjects.ID"},
			{FieldName: "confirmableEmailAddress", DataType: "valueobjects.ConfirmableEmailAddress"},
			{FieldName: "personName", DataType: "valueobjects.PersonName"},
		},
	},
	{
		EventType:    "emailAddressConfirmed",
		EventFactory: "EmailAddressWasConfirmed",
		Fields: []Field{
			{FieldName: "id", DataType: "valueobjects.ID"},
			{FieldName: "emailAddress", DataType: "valueobjects.EmailAddress"},
		},
	},
}

type Config struct {
	RelativeOutputPath string
	AggregateFactory   string
	Events             []Event
}

var config = Config{
	RelativeOutputPath: "../../domain",
	AggregateFactory:   "NewUnregisteredCustomer",
}

func main() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	for _, event := range events {
		event.AggregateFactory = config.AggregateFactory

		t := template.New(event.EventType)
		t = t.Funcs(
			template.FuncMap{
				"title":      strings.Title,
				"methodName": methodName,
				"eventName":  t.Name,
			},
		)

		t, err = t.Parse(tmpl)
		die(err)

		outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + ".go")
		die(err)

		// comment for generating without formating
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

		// uncomment for generating without formating
		//err = t.Execute(outFile, event)
		//die(err)

		// uncomment for testing, outputs to stdout
		//err = t.Execute(os.Stdout, event)
		//die(err)
	}

}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var tmpl = `
package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

type {{title eventName}} interface {
	{{range .Fields}}{{methodName .DataType}}() {{.DataType}}
	{{end}}
	shared.DomainEvent
}

type {{eventName}} struct {
	{{range .Fields}}{{.FieldName}} {{.DataType}}
	{{end}}
	meta *shared.DomainEventMeta
}

func {{.EventFactory}}(
	{{range .Fields}} {{.FieldName}} {{.DataType}},
	{{end -}}
) *{{eventName}} {

	{{eventName}} := &{{eventName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}},
		{{end -}}
	}

	{{eventName}}.meta = shared.NewDomainEventMeta(id, {{.AggregateFactory}}(), {{eventName}})

	return {{eventName}}
}

{{range .Fields}}
func ({{eventName}} *{{eventName}}) {{methodName .DataType}}() {{.DataType}} {
	return {{eventName}}.{{.FieldName}}
}
{{end}}

func ({{eventName}} *{{eventName}}) Identifier() shared.AggregateIdentifier {
	return {{eventName}}.meta.Identifier
}

func ({{eventName}} *{{eventName}}) EventName() string {
	return {{eventName}}.meta.EventName
}

func ({{eventName}} *{{eventName}}) OccurredAt() string {
	return {{eventName}}.meta.OccurredAt
}
`
