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
		EventType:    "Registered",
		EventFactory: "ItWasRegistered",
		Fields: []Field{
			{FieldName: "id", DataType: "*valueobjects.ID"},
			{FieldName: "confirmableEmailAddress", DataType: "*valueobjects.ConfirmableEmailAddress"},
			{FieldName: "personName", DataType: "*valueobjects.PersonName"},
		},
	},
	{
		EventType:    "EmailAddressConfirmed",
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
	RelativeOutputPath: "../../domain/events",
}

func main() {
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
				"lcFirst":    lcFirst,
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

const {{lcFirst eventName}}AggregateName = "Customer"

type {{eventName}} struct {
	{{range .Fields}}{{.FieldName}} {{.DataType}}
	{{end}}
	meta *shared.DomainEventMeta
}

/*** Factory Methods ***/

func {{.EventFactory}}(
	{{range .Fields}}{{.FieldName}} {{.DataType}},
	{{end -}}
) *{{eventName}} {

	{{lcFirst eventName}} := &{{eventName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}},
		{{end -}}
	}

	{{lcFirst eventName}}.meta = shared.NewDomainEventMeta(
		id.String(),
		{{lcFirst eventName}},
		{{lcFirst eventName}}AggregateName,
	)

	return {{lcFirst eventName}}
}

/*** Getter Methods ***/

{{range .Fields}}
func ({{lcFirst eventName}} *{{eventName}}) {{methodName .DataType}}() {{.DataType}} {
	return {{lcFirst eventName}}.{{.FieldName}}
}
{{end}}

/*** Implement shared.DomainEvent ***/

func ({{lcFirst eventName}} *{{eventName}}) Identifier() string {
	return {{lcFirst eventName}}.meta.Identifier
}

func ({{lcFirst eventName}} *{{eventName}}) EventName() string {
	return {{lcFirst eventName}}.meta.EventName
}

func ({{lcFirst eventName}} *{{eventName}}) OccurredAt() string {
	return {{lcFirst eventName}}.meta.OccurredAt
}

/*** Implement json.Marshaler ***/

func ({{lcFirst eventName}} *{{eventName}}) MarshalJSON() ([]byte, error) {
	data := &struct {
		{{range .Fields}}{{methodName .DataType}} {{.DataType}} {{tick}}json:"{{.FieldName}}"{{tick}}
		{{end -}}
		Meta *shared.DomainEventMeta {{tick}}json:"meta"{{tick}}
	}{
		{{range .Fields}}{{methodName .DataType}}: {{lcFirst eventName}}.{{.FieldName}},
		{{end -}}
		Meta: {{lcFirst eventName}}.meta,
	}

	return json.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func ({{lcFirst eventName}} *{{eventName}}) UnmarshalJSON(data []byte) error {
	values := &struct {
		{{range .Fields}}{{methodName .DataType}} {{.DataType}} {{tick}}json:"{{.FieldName}}"{{tick}}
		{{end -}}
		Meta *shared.DomainEventMeta {{tick}}json:"meta"{{tick}}
	}{}

	if err := json.Unmarshal(data, values); err != nil {
		return err
	}

	{{range .Fields}}{{lcFirst eventName}}.{{.FieldName}} = values.{{methodName .DataType}}
	{{end -}}

	{{lcFirst eventName}}.meta = values.Meta

	return nil
}
`
