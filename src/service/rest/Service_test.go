package rest_test

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	customergrpc "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/src/service"
	grpcService "github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/service/rest"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/connectivity"
)

func TestStartRestServer(t *testing.T) {
	mockedExistingCustomerID := "11111111"

	logger := shared.NewNilLogger()
	config := service.MustBuildConfigFromEnv(logger)

	exitWasCalled := false
	exitFn := func() {
		exitWasCalled = true
	}
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)

	runGRPCServer(config, logger, mockedExistingCustomerID)

	grpcClientConn := rest.MustDialGRPCContext(config, logger, ctx, cancelFn)

	terminateDelay := time.Millisecond * 100

	s := rest.InitService(config, logger, exitFn, ctx, cancelFn, grpcClientConn)

	Convey("Start the REST server as a goroutine", t, func() {
		go s.StartRestServer()

		Convey("REST server should handle successful requests serving a static file", func() {
			var err error
			var resp *resty.Response

			hostAndPort := config.REST.HostAndPort

			client := resty.New()

			resp, err = client.R().
				Get(fmt.Sprintf("http://%s/v1/customer/swagger.json", config.REST.HostAndPort))
			So(err, ShouldBeNil)
			So(resp.StatusCode(), ShouldEqual, 200)

			Convey("REST server should handle successful requests served via gRPC", func() {
				resp, err = client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(`{"emailAddress": "anton+10@stoeckl.de", "familyName": "Stöckl", "givenName": "Anton"}`).
					Post(fmt.Sprintf("http://%s/v1/customer", hostAndPort))

				So(err, ShouldBeNil)
				So(resp.StatusCode(), ShouldEqual, 200)

				resp, err = client.R().
					Delete(fmt.Sprintf("http://%s/v1/customer/%s", hostAndPort, mockedExistingCustomerID))
				So(err, ShouldBeNil)
				So(resp.StatusCode(), ShouldEqual, 200)

				Convey("REST server should handle failed requests served via gRPC", func() {
					notExistingCustomerID := "66666666"

					resp, _ := client.R().
						Get(fmt.Sprintf("http://%s/v1/customer/%s", hostAndPort, notExistingCustomerID))
					So(resp.StatusCode(), ShouldEqual, 404)

					Convey(fmt.Sprintf("It should wait for stop signal (scheduled after %s)", terminateDelay), func() {
						start := time.Now()
						go func() {
							time.Sleep(terminateDelay)
							_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
						}()

						s.WaitForStopSignal()

						So(time.Now(), ShouldNotHappenWithin, terminateDelay, start)

						Convey("Stop signal should issue Shutdown", func() {
							Convey("Shutdown should cancel the context", func() {
								So(errors.Is(ctx.Err(), context.Canceled), ShouldBeTrue)

								Convey("Shutdown should stop REST server", func() {
									resp, err := client.R().
										Get(fmt.Sprintf("http://%s/v1/customer/1234", hostAndPort))

									So(err, ShouldBeError)
									So(err.Error(), ShouldContainSubstring, "connection refused")
									So(resp.StatusCode(), ShouldBeZeroValue)

									Convey("Shutdown should close the grpc client connection", func() {
										So(grpcClientConn.GetState(), ShouldResemble, connectivity.Shutdown)

										Convey("Shutdown should call exit", func() {
											So(exitWasCalled, ShouldBeTrue)
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
}

/*** Helper functions ***/

func runGRPCServer(config *service.Config, logger *shared.Logger, mockedExistingCustomerID string) {
	diContainer := service.MustBuildDIContainer(
		config,
		logger,
		service.ReplaceGRPCCustomerServer(grpcCustomerServerStub(mockedExistingCustomerID)),
	)
	grpcSvc := grpcService.InitService(config, logger, func() {}, diContainer)
	go grpcSvc.StartGRPCServer()
}

func grpcCustomerServerStub(mockedExistingCustomerID string) customergrpc.CustomerServer {
	customerServer := customergrpc.NewCustomerServer(
		func(emailAddress, givenName, familyName string) (value.CustomerID, error) {
			return value.GenerateCustomerID(), nil
		},
		func(customerID, confirmationHash string) error {
			return nil
		},
		func(customerID, emailAddress string) error {
			return nil
		},
		func(customerID, givenName, familyName string) error {
			return nil
		},
		func(customerID string) error {
			return nil
		},
		func(customerID string) (customer.View, error) {
			switch customerID {
			case mockedExistingCustomerID:
				return customer.View{ID: customerID}, nil
			default:
				return customer.View{}, shared.ErrNotFound
			}
		},
	)

	return customerServer
}
