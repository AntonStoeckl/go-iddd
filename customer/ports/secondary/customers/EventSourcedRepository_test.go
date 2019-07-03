package customers_test

import (
	"errors"
	"fmt"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/customer/ports/secondary/customers"
	"go-iddd/customer/ports/secondary/customers/mocks"
	"go-iddd/shared"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestEventSourcedRepositoryRegister(t *testing.T) {
	Convey("When a new Customer is registered", t, func() {
		id := values.GenerateID()

		customer, err := buildRegisteredCustomerWith(id)
		So(err, ShouldBeNil)

		eventStore := shared.NewInMemoryEventStore("customer")

		persistableCustomers := customers.NewEventSourcedRepository(
			eventStore,
			domain.ReconstituteCustomerFrom,
		)

		Convey("And when appending the recorded events succeed", func() {
			err := persistableCustomers.Register(customer)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And when the same Customer is registered again", func() {
					err := persistableCustomers.Register(customer)

					Convey("It should fail", func() {
						So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
					})
				})
			})
		})

		Convey("When appending the recorded events fails", func() {
			Convey("And when if fails with a concurrency conflict", func() {
				eventStore.FailOnceWith(shared.ErrConcurrencyConflict)
				err := persistableCustomers.Register(customer)

				Convey("It should report a duplicate Customer", func() {
					So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
				})
			})

			Convey("And when it fails with a technical problem", func() {
				expectedErr := errors.New("mocked error")
				eventStore.FailOnceWith(expectedErr)
				err := persistableCustomers.Register(customer)

				Convey("It should report the technical problem", func() {
					So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				})
			})
		})
	})
}

func TestEventSourcedRepositoryOf(t *testing.T) {
	Convey("Given an existing Customer", t, func() {
		id := values.GenerateID()

		customer, err := buildRegisteredCustomerWith(id)
		So(err, ShouldBeNil)

		eventStore := shared.NewInMemoryEventStore("customer")

		persistableCustomers := customers.NewEventSourcedRepository(eventStore, domain.ReconstituteCustomerFrom)

		err = persistableCustomers.Register(customer)
		So(err, ShouldBeNil)

		Convey("And given the Customer is not cached", func() {
			// recreate the repository to reset the cache (cache)
			persistableCustomers = customers.NewEventSourcedRepository(eventStore, domain.ReconstituteCustomerFrom)

			Convey("When the Customer is retrieved", func() {
				Convey("And when loading the event stream succeeds", func() {
					Convey("And when the customer can be reconstituted", func() {
						customer, err := persistableCustomers.Of(id)

						Convey("It should succeed", func() {
							So(err, ShouldBeNil)
							So(customer.AggregateID().Equals(id), ShouldBeTrue)
							So(customer.StreamVersion(), ShouldEqual, uint(1))
						})
					})

					Convey("And when the customer can not be reconstituted", func() {
						expectedErr := errors.New("mocked error")
						customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
							return nil, expectedErr
						}

						persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

						customer, err := persistableCustomers.Of(id)

						Convey("It should fail", func() {
							So(xerrors.Is(err, expectedErr), ShouldBeTrue)
							So(customer, ShouldBeNil)
						})
					})
				})

				Convey("And when loading the event stream fails", func() {
					expectedErr := errors.New("mocked error")
					eventStore.FailOnceWith(expectedErr)
					customer, err := persistableCustomers.Of(id)

					Convey("It should fail", func() {
						So(xerrors.Is(err, expectedErr), ShouldBeTrue)
						So(customer, ShouldBeNil)
					})
				})
			})
		})

		Convey("And given the Customer is cached", func() {
			Convey("And given the Customer was concurrently modified", func() {
				otherRepo := customers.NewEventSourcedRepository(
					eventStore,
					domain.ReconstituteCustomerFrom,
				)

				sameCustomer, err := otherRepo.Of(id)
				So(err, ShouldBeNil)

				newEmailAddress := fmt.Sprintf("john+changed+%s@doe.com", id.String())
				changeEmailAddress, err := commands.NewChangeEmailAddress(id.String(), newEmailAddress)

				err = sameCustomer.Execute(changeEmailAddress)
				So(err, ShouldBeNil)

				err = otherRepo.Persist(sameCustomer)
				So(err, ShouldBeNil)

				Convey("When the Customer is retrieved", func() {
					Convey("And when loading the new events succeeds", func() {
						customer, err := persistableCustomers.Of(id)

						Convey("It should retrieve the up-to-date Customer", func() {
							So(err, ShouldBeNil)
							So(customer.AggregateID().Equals(id), ShouldBeTrue)
							So(customer.StreamVersion(), ShouldEqual, uint(2))
						})
					})

					Convey("And when loading the new events fails", func() {
						expectedErr := errors.New("mocked error")
						eventStore.FailOnceWith(expectedErr)
						customer, err := persistableCustomers.Of(id)

						Convey("It should fail", func() {
							So(xerrors.Is(err, expectedErr), ShouldBeTrue)
							So(customer, ShouldBeNil)
						})
					})
				})
			})
		})
	})

	Convey("Given a not existing Customer", t, func() {
		id := values.GenerateID()

		eventStore := new(mocks.EventStore)
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, domain.ReconstituteCustomerFrom)

		Convey("When the event stream is loaded", func() {
			eventStream := shared.DomainEvents{}
			eventStore.On("LoadEventStream", id).Return(eventStream, nil).Once()
			customer, err := persistableCustomers.Of(id)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)
				So(customer, ShouldBeNil)
			})
		})
	})
}

func TestEventSourcedRepositoryPersist(t *testing.T) {
	Convey("When a Customer is persisted", t, func() {
		id := values.GenerateID()

		customer, err := buildRegisteredCustomerWith(id)
		So(err, ShouldBeNil)

		registered := customer.RecordedEvents(false)[0].(*events.Registered)

		eventStore := new(mocks.EventStore)
		eventStore.
			On("AppendToStream", customer.RecordedEvents(false)).
			Return(nil).
			Once()

		persistableCustomers := customers.NewEventSourcedRepository(eventStore, domain.ReconstituteCustomerFrom)

		err = persistableCustomers.Register(customer)
		So(err, ShouldBeNil)

		emailAddressConfirmed, err := commands.NewConfirmEmailAddress(
			id.String(),
			registered.ConfirmableEmailAddress().EmailAddress(),
			registered.ConfirmableEmailAddress().ConfirmationHash(),
		)
		So(err, ShouldBeNil)

		err = customer.Execute(emailAddressConfirmed)
		So(err, ShouldBeNil)

		Convey("And when appending the recorded events succeed", func() {
			eventStore.
				On("AppendToStream", customer.RecordedEvents(false)).
				Return(nil).
				Once()
			err := persistableCustomers.Persist(customer)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)
			})
		})

		Convey("And when appending the recorded events fails", func() {
			expectedErr := errors.New("mocked error")
			eventStore.
				On("AppendToStream", customer.RecordedEvents(false)).
				Return(expectedErr).
				Once()

			err := persistableCustomers.Persist(customer)

			Convey("It should fail", func() {
				So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}

func TestEventSourcedRepositoryCacheRaceConditions(t *testing.T) {
	Convey("When the cache has concurrent read/write requests", t, func() {
		persistableCustomers := customers.NewEventSourcedRepository(
			shared.NewInMemoryEventStore("customer"),
			domain.ReconstituteCustomerFrom,
		)

		Convey("It should not report race conditions", func() {
			wg := new(sync.WaitGroup)

			concurrent := func(wg *sync.WaitGroup) {
				defer wg.Done()

				id := values.GenerateID()
				customer, _ := buildRegisteredCustomerWith(id)

				_ = persistableCustomers.Register(customer)
				_, _ = persistableCustomers.Of(id)
			}

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go concurrent(wg)
			}

			wg.Wait()
		})
	})
}

func buildRegisteredCustomerWith(id *values.ID) (domain.Customer, error) {
	emailAddress := fmt.Sprintf("john+%s@doe.com", id.String())
	givenName := "John"
	familyName := "Doe"

	register, err := commands.NewRegister(id.String(), emailAddress, givenName, familyName)
	if err != nil {
		return nil, err
	}

	customer := domain.NewCustomerWith(register)

	return customer, nil
}
