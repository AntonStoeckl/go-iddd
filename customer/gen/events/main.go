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

const mode = "format" // "format", "noformat", "stoud"

type Field struct {
	FieldName string
	DataType  string
}

type Event struct {
	EventType    string
	EventFactory string
	Fields       []Field
}

var events = []Event{
	{
		EventType:    "registered",
		EventFactory: "ItWasRegistered",
		Fields: []Field{
			{FieldName: "id", DataType: "*valueobjects.ID"},
			{FieldName: "confirmableEmailAddress", DataType: "*valueobjects.ConfirmableEmailAddress"},
			{FieldName: "personName", DataType: "*valueobjects.PersonName"},
		},
	},
	{
		EventType:    "emailAddressConfirmed",
		EventFactory: "EmailAddressWasConfirmed",
		Fields: []Field{
			{FieldName: "id", DataType: "*valueobjects.ID"},
			{FieldName: "emailAddress", DataType: "*valueobjects.EmailAddress"},
		},
	},
}

type Config struct {
	RelativeOutputPath string
	Events             []Event
}

var config = Config{
	RelativeOutputPath: "../../domain",
}

func main() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	tick := func() string {
		return "`"
	}

	for _, event := range events {
		t := template.New(event.EventType)
		t = t.Funcs(
			template.FuncMap{
				"title":      strings.Title,
				"methodName": methodName,
				"eventName":  t.Name,
				"tick":       tick,
			},
		)

		t, err = t.Parse(eventTemplate)
		die(err)

		switch mode {
		case "format":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + ".go")
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
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + ".go")
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

var eventTemplate = `
package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

const {{eventName}}AggregateName = "Customer"

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
	{{range .Fields}}{{.FieldName}} {{.DataType}},
	{{end -}}
) *{{eventName}} {

	{{eventName}} := &{{eventName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}},
		{{end -}}
	}

	{{eventName}}.meta = shared.NewDomainEventMeta(
		id.String(),
		{{eventName}},
		{{eventName}}AggregateName,
	)

	return {{eventName}}
}

{{range .Fields}}
func ({{eventName}} *{{eventName}}) {{methodName .DataType}}() {{.DataType}} {
	return {{eventName}}.{{.FieldName}}
}
{{end}}

func ({{eventName}} *{{eventName}}) Identifier() string {
	return {{eventName}}.meta.Identifier
}

func ({{eventName}} *{{eventName}}) EventName() string {
	return {{eventName}}.meta.EventName
}

func ({{eventName}} *{{eventName}}) OccurredAt() string {
	return {{eventName}}.meta.OccurredAt
}

func ({{eventName}} *{{eventName}}) MarshalJSON() ([]byte, error) {
	data := &struct {
		{{range .Fields}}{{methodName .DataType}} {{.DataType}} {{tick}}json:"{{.FieldName}}"{{tick}}
		{{end -}}
		Meta *shared.DomainEventMeta {{tick}}json:"meta"{{tick}}
	}{
		{{range .Fields}}{{methodName .DataType}}: {{eventName}}.{{.FieldName}},
		{{end -}}
		Meta: {{eventName}}.meta,
	}

	return json.Marshal(data)
}

func Unmarshal{{title eventName}}FromJSON(jsonData []byte) ({{title eventName}}, error) {
	var err error
	var data map[string]interface{}

	{{range .Fields}}var {{.FieldName}} {{.DataType}}
	{{end -}}
	var meta *shared.DomainEventMeta

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	for key, value := range data {
		switch key {
		{{range .Fields}}case "{{.FieldName}}":
			if {{.FieldName}}, err = valueobjects.Unmarshal{{methodName .DataType}}(value); err != nil {
				return nil, err
			}
		{{end -}}
		case "meta":
			if meta, err = shared.UnmarshalDomainEventMeta(value); err != nil {
				return nil, err
			}
		}
	}

	{{eventName}} := &{{eventName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}},
		{{end -}}
		meta: meta,
	}

	return {{eventName}}, nil
}
`
