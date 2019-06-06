package customers_test

import (
	"errors"
	"go-iddd/customer/domain"
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
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)
		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion)
		currentStreamVersion++
		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress, currentStreamVersion)
		recordedEvents := shared.DomainEvents{registered, emailAddressConfirmed}

		customer := new(mocks.Customer)
		customer.On("AggregateID").Return(id)
		customer.On("RecordedEvents").Return(recordedEvents)

		eventStore := new(mocks.EventStore)
		customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
			return customer, nil
		}
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

		Convey("And when appending the recorded events succeed", func() {
			eventStore.On("AppendToStream", id, recordedEvents).Return(nil).Once()
			err := persistableCustomers.Register(customer)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)

				Convey("And when the same Customer is registered again", func() {
					// eventStore.AppendToStream() is mocked to be called only Once, so it would fail if it's called again
					err := persistableCustomers.Register(customer)

					Convey("It should fail", func() {
						So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
						So(eventStore.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})

			Convey("And when the Customer is retrieved", func() {
				// eventStore.LoadEventStream() is not mocked, so it would fail if it's called
				customer, err := persistableCustomers.Of(id)

				Convey("It should be retrieved from the identity map (cache)", func() {
					So(err, ShouldBeNil)
					So(eventStore.AssertExpectations(t), ShouldBeTrue)
					So(customer.AggregateID().Equals(id), ShouldBeTrue)
				})
			})
		})

		Convey("When appending the recorded events fails", func() {
			Convey("And when if fails with a concurrency conflict", func() {
				eventStore.On("AppendToStream", id, recordedEvents).Return(shared.ErrConcurrencyConflict).Once()
				err := persistableCustomers.Register(customer)

				Convey("It should report a duplicate Customer", func() {
					So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
					So(eventStore.AssertExpectations(t), ShouldBeTrue)
				})
			})

			Convey("And when it fails with a technical problem", func() {
				expectedErr := errors.New("mocked error")
				eventStore.On("AppendToStream", id, recordedEvents).Return(expectedErr).Once()
				err := persistableCustomers.Register(customer)

				Convey("It should report the technical problem", func() {
					So(xerrors.Is(err, expectedErr), ShouldBeTrue)
					So(eventStore.AssertExpectations(t), ShouldBeTrue)
				})
			})
		})
	})
}

func TestEventSourcedRepositoryOf(t *testing.T) {
	Convey("Given an existing Customer", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)
		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion)
		currentStreamVersion++
		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress, currentStreamVersion)
		recordedEvents := shared.DomainEvents{registered, emailAddressConfirmed}

		customer := new(mocks.Customer)
		customer.On("AggregateID").Return(id)
		customer.On("RecordedEvents").Return(recordedEvents)

		eventStore := new(mocks.EventStore)
		customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
			return customer, nil
		}
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

		eventStore.On("AppendToStream", id, recordedEvents).Return(nil).Once()
		err = persistableCustomers.Register(customer)
		So(err, ShouldBeNil)

		// recreate the repository to reset the identityMap (cache)
		persistableCustomers = customers.NewEventSourcedRepository(eventStore, customerFactory)

		Convey("When the Customer is retrieved", func() {
			Convey("And when loading the event stream succeeds", func() {
				Convey("And when the customer can be reconstituted", func() {
					eventStore.On("LoadEventStream", id).Return(recordedEvents, nil).Once()
					customer, err := persistableCustomers.Of(id)

					Convey("It should succeed", func() {
						So(err, ShouldBeNil)
						So(eventStore.AssertExpectations(t), ShouldBeTrue)
						So(customer.AggregateID().Equals(id), ShouldBeTrue)
					})

					Convey("And when the Customer is retrieved again", func() {
						// eventStore.LoadEventStream() is mocked to be called only Once, so it would fail if it's called again
						customer, err := persistableCustomers.Of(id)

						Convey("It should be retrieved from the identity map (cache)", func() {
							So(err, ShouldBeNil)
							So(eventStore.AssertExpectations(t), ShouldBeTrue)
							So(customer.AggregateID().Equals(id), ShouldBeTrue)
						})
					})
				})

				Convey("And when the customer can not be reconstituted", func() {
					eventStore.On("LoadEventStream", id).Return(recordedEvents, nil).Once()
					expectedErr := errors.New("mocked error")
					customerFactory = func(eventStream shared.DomainEvents) (domain.Customer, error) {
						return nil, expectedErr
					}

					persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

					customer, err := persistableCustomers.Of(id)

					Convey("It should fail", func() {
						So(xerrors.Is(err, expectedErr), ShouldBeTrue)
						So(eventStore.AssertExpectations(t), ShouldBeTrue)
						So(customer, ShouldBeNil)
					})
				})
			})

			Convey("And when loading the event stream fails", func() {
				expectedErr := errors.New("mocked error")
				eventStore.On("LoadEventStream", id).Return(nil, expectedErr).Once()
				customer, err := persistableCustomers.Of(id)

				Convey("It should fail", func() {
					So(xerrors.Is(err, expectedErr), ShouldBeTrue)
					So(eventStore.AssertExpectations(t), ShouldBeTrue)
					So(customer, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given a not existing Customer", t, func() {
		id := values.GenerateID()
		customer := new(mocks.Customer)

		eventStore := new(mocks.EventStore)
		customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
			return customer, nil
		}
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

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
		recordedEvents := shared.DomainEvents{}

		customer := new(mocks.Customer)
		customer.On("AggregateID").Return(id)
		customer.On("RecordedEvents").Return(recordedEvents)

		eventStore := new(mocks.EventStore)
		customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
			return customer, nil
		}
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

		Convey("And when appending the recorded events succeed", func() {
			eventStore.On("AppendToStream", id, recordedEvents).Return(nil).Once()
			err := persistableCustomers.Persist(customer)

			Convey("It should register the customer", func() {
				So(err, ShouldBeNil)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)
			})
		})

		Convey("And when appending the recorded events fails", func() {
			eventStream := shared.DomainEvents{}
			expectedErr := errors.New("mocked error")
			eventStore.On("AppendToStream", id, eventStream).Return(expectedErr).Once()
			err := persistableCustomers.Persist(customer)

			Convey("It should fail", func() {
				So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				So(eventStore.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}

func TestEventSourcedRepositoryIdentityMapRaceConditions(t *testing.T) {
	Convey("When the identityMap has concurrent read/write requests", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.NewPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)
		registered := events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion)
		currentStreamVersion++
		emailAddressConfirmed := events.EmailAddressWasConfirmed(id, emailAddress, currentStreamVersion)
		recordedEvents := shared.DomainEvents{registered, emailAddressConfirmed}

		customer := new(mocks.Customer)
		customer.On("AggregateID").Return(id)
		customer.On("RecordedEvents").Return(recordedEvents)

		eventStore := new(mocks.EventStore)
		eventStore.On("AppendToStream", id, recordedEvents).Return(nil)
		eventStore.On("LoadEventStream", id).Return(recordedEvents, nil)
		customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
			return customer, nil
		}
		persistableCustomers := customers.NewEventSourcedRepository(eventStore, customerFactory)

		Convey("It should not report race conditions", func() {
			wg := new(sync.WaitGroup)

			concurrent := func(wg *sync.WaitGroup) {
				defer wg.Done()

				_ = persistableCustomers.Register(customer)
				_, _ = persistableCustomers.Of(id)
			}

			for i := 0; i < 2; i++ {
				wg.Add(1)
				go concurrent(wg)
			}

			wg.Wait()
		})
	})
}
