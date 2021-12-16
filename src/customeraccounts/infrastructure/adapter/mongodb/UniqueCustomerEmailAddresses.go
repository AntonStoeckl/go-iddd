package mongodb

import (
	"context"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/cockroachdb/errors"
)

type UniqueCustomerEmailAddresses struct {
	uniqueEmailAddressesCollection    *mongo.Collection
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
}

func NewUniqueCustomerEmailAddresses(
	uniqueEmailAddressesCollection *mongo.Collection,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
) *UniqueCustomerEmailAddresses {
	// @TODO as init script
	unigeEmailIndex := mongo.IndexModel{
		Keys: primitive.M{
			"email_address": 1, // index in descending order
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}
	if _, e := uniqueEmailAddressesCollection.Indexes().CreateOne(context.Background(), unigeEmailIndex); e != nil {
		panic(e)
	}
	return &UniqueCustomerEmailAddresses{
		uniqueEmailAddressesCollection:    uniqueEmailAddressesCollection,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}
}

func (s *UniqueCustomerEmailAddresses) AssertUniqueEmailAddress(recordedEvents []es.DomainEvent, tx mongo.SessionContext) error {
	wrapWithMsg := "assertUniqueEmailAddresse"

	assertions := s.buildUniqueEmailAddressAssertions(recordedEvents...)

	for _, assertion := range assertions {
		switch assertion.DesiredAction() {
		case customer.ShouldAddUniqueEmailAddress:
			if err := s.tryToAdd(assertion.EmailAddressToAdd(), assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case customer.ShouldReplaceUniqueEmailAddress:
			if err := s.tryToReplace(assertion.EmailAddressToAdd(), assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case customer.ShouldRemoveUniqueEmailAddress:
			if err := s.remove(assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) PurgeUniqueEmailAddress(customerID value.CustomerID, tx mongo.SessionContext) error {
	return s.remove(customerID, tx)
}

func (s *UniqueCustomerEmailAddresses) tryToAdd(
	emailAddress value.UnconfirmedEmailAddress,
	customerID value.CustomerID,
	tx mongo.SessionContext,
) error {

	// queryTemplate := `INSERT INTO %tablename% VALUES ($1, $2)`
	// query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesCollection, 1)

	_, err := s.uniqueEmailAddressesCollection.InsertOne(tx,
		primitive.D{
			primitive.E{Key: "email_address", Value: emailAddress.String()},
			primitive.E{Key: "customer_id", Value: customerID.String()},
		})

	if err != nil {
		return s.mapUniqueEmailAddressMongodbErrors(err)
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) tryToReplace(
	emailAddress value.UnconfirmedEmailAddress,
	customerID value.CustomerID,
	tx mongo.SessionContext,
) error {

	// queryTemplate := `UPDATE %tablename% set email_address = $1 where customer_id = $2`
	// query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesCollection, 1)

	_, err := s.uniqueEmailAddressesCollection.UpdateOne(tx,
		primitive.D{
			primitive.E{Key: "customer_id", Value: customerID.String()},
		},
		primitive.D{
			primitive.E{Key: "$set", Value: primitive.D{
				primitive.E{Key: "email_address", Value: emailAddress.String()},
			}},
		})

	if err != nil {
		return s.mapUniqueEmailAddressMongodbErrors(err)
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) remove(
	customerID value.CustomerID,
	tx mongo.SessionContext,
) error {

	// queryTemplate := `DELETE FROM %tablename% where customer_id = $1`
	// query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesCollection, 1)

	_, err := s.uniqueEmailAddressesCollection.DeleteOne(tx, primitive.D{
		primitive.E{Key: "customer_id", Value: customerID.String()}},
	)

	if err != nil {
		return s.mapUniqueEmailAddressMongodbErrors(err)
	}

	return nil
}
func (s *UniqueCustomerEmailAddresses) mapUniqueEmailAddressMongodbErrors(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(mongo.WriteException); ok {
		if actualErr.HasErrorCode(11000) {
			return errors.Mark(errors.New("duplicate email address"), shared.ErrDuplicate)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
