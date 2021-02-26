package customeraccounts_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	. "github.com/smartystreets/goconvey/convey"
)

func TestScenarios_ForRegisteringIdentities(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		h := application.NewIdentityCommandHandler()
		id := value.GenerateIdentityID()
		ea := "kevin@ball.net"
		pw := "superSecretPW123"

		Convey("\nSCENARIO: A prospective Customer registers his identity", func() {
			err := h.RegisterIdentity(id, ea, pw)
			So(err, ShouldBeNil)
		})
	})
}
