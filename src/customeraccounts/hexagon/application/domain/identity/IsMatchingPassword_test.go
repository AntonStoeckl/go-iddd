package identity_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMatchesIdentityPassword(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		identityID := value.GenerateIdentityID()
		emailAddress := "kevin@ball.com"
		plainPassword := "superSecretPW123"
		differentPlainPassword := "differentSecretPW"

		query, err := domain.BuildIsMatchingPasswordForIdentity(emailAddress, plainPassword)
		So(err, ShouldBeNil)

		queryWithDifferentPassword, err := domain.BuildIsMatchingPasswordForIdentity(emailAddress, differentPlainPassword)
		So(err, ShouldBeNil)

		emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
		So(err, ShouldBeNil)

		hashedPassword, err := value.HashedPasswordFromPlainPassword(value.RebuildPlainPassword(plainPassword))
		So(err, ShouldBeNil)

		identityRegistered := domain.BuildIdentityRegistered(
			identityID,
			emailAddressValue,
			hashedPassword,
			es.GenerateMessageID(),
			1,
		)

		identityDeleted := domain.BuildIdentityDeleted(
			identityID,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO: Match an Identity's password with a matching supplied password", func() {
			Convey("Given IdentityRegistered", func() {
				eventStream := es.EventStream{identityRegistered}

				Convey("When IsMatchingPassword with matching password", func() {
					err := identity.IsMatchingPassword(eventStream, query)

					Convey("Then it should match", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("\nSCENARIO: Match an Identity's password with a different supplied password", func() {
			Convey("Given IdentityRegistered", func() {
				eventStream := es.EventStream{identityRegistered}

				Convey("When IsMatchingPassword with different password", func() {
					err := identity.IsMatchingPassword(eventStream, queryWithDifferentPassword)

					Convey("Then it should not match", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrInvalidCredentials), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Match an Identity's password when the Identity was deleted", func() {
			Convey("Given IdentityRegistered", func() {
				eventStream := es.EventStream{identityRegistered}

				Convey("and IdentityDeleted", func() {
					eventStream = append(eventStream, identityDeleted)

					Convey("When IsMatchingPassword with matching password", func() {
						err := identity.IsMatchingPassword(eventStream, query)

						Convey("Then it should fail", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrInvalidCredentials), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
