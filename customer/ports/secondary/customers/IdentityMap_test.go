package customers

import (
	"fmt"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIdentityMap_Memoize_And_MemoizedCustomerOf(t *testing.T) {
	Convey("Given a Customer was not stored to the IdentityMap", t, func() {
		identityMap := NewIdentityMap()
		id := values.GenerateID()
		register, err := commands.NewRegister(id.String(), fmt.Sprintf("john+%s@doe.com", id.String()), "John", "Doe")
		So(err, ShouldBeNil)
		customer := domain.NewCustomerWith(register)

		Convey("The Customer should not be found in the IdentityMap", func() {
			actualCustomer, found := identityMap.MemoizedCustomerOf(id)
			So(found, ShouldBeFalse)
			So(actualCustomer, ShouldBeNil)
		})

		Convey("When the Customer is momoized", func() {
			identityMap.Memoize(customer)

			Convey("The Customer should be found in the IdentityMap", func() {
				actualCustomer, found := identityMap.MemoizedCustomerOf(id)
				So(found, ShouldBeTrue)
				So(actualCustomer, ShouldResemble, customer)
			})
		})
	})
}

func TestIdentityMap_RaceConditions(t *testing.T) {
	Convey("When the IdentityMap has concurrent read/write requests", t, func() {
		identityMap := NewIdentityMap()

		Convey("It should not have race conditions", func() {
			wg := new(sync.WaitGroup)

			concurrent := func(wg *sync.WaitGroup) {
				defer wg.Done()

				id := values.GenerateID()
				register, err := commands.NewRegister(id.String(), fmt.Sprintf("john+%s@doe.com", id.String()), "John", "Doe")
				if err != nil {
					panic("creating a register command for test failed")
				}
				customer := domain.NewCustomerWith(register)

				identityMap.Memoize(customer)
				_, _ = identityMap.MemoizedCustomerOf(id)
			}

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go concurrent(wg)
			}

			wg.Wait()
		})
	})
}
