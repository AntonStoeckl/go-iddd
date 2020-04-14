package cmd

import (
	"database/sql"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDIContainer(t *testing.T) {
	Convey("When a DIContainer is created with a nil postgresDBConn", t, func() {
		var conn *sql.DB

		marshalDomainEvent := func(event es.DomainEvent) ([]byte, error) {
			return nil, nil
		}

		unmarshalDomainEvent := func(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) {
			return nil, nil
		}

		_, err := NewDIContainer(conn, marshalDomainEvent, unmarshalDomainEvent)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
		})
	})
}
