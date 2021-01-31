package rest_test

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/service"
	"github.com/AntonStoeckl/go-iddd/src/service/rest"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/connectivity"
)

func TestStartRestServer(t *testing.T) {
	logger := shared.NewNilLogger()
	config := service.MustBuildConfigFromEnv(logger)

	exitWasCalled := false
	exitFn := func() {
		exitWasCalled = true
	}
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)

	grpcClientConn := rest.MustDialGRPCContext(config, logger, ctx, cancelFn)

	terminateDelay := time.Millisecond * 100

	s := rest.InitService(config, logger, exitFn, ctx, cancelFn, grpcClientConn)

	Convey("Start the REST server as a goroutine", t, func() {
		go s.StartRestServer()

		Convey("REST server should handle successful requests", func() {
			var err error
			var resp *resty.Response

			hostAndPort := config.REST.HostAndPort
			mockedCustomerID := "11111111"
			notExistingCustomerID := "66666666"

			client := resty.New()

			// TODO: why the hack does this 404?
			//resp, err = client.R().
			//	Get(fmt.Sprintf("http://%s/v1/customer/swagger.json", config.REST.HostAndPort))
			//So(err, ShouldBeNil)
			//So(resp.StatusCode(), ShouldEqual, 200)

			result := &struct {
				ID string `json:"id"`
			}{ID: mockedCustomerID}

			resp, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(`{"emailAddress": "anton+10@stoeckl.de", "familyName": "Stöckl", "givenName": "Anton"}`).
				SetResult(result).
				Post(fmt.Sprintf("http://%s/v1/customer", hostAndPort))

			So(err, ShouldBeNil)
			So(resp.StatusCode(), ShouldEqual, 200)

			resp, err = client.R().
				Delete(fmt.Sprintf("http://%s/v1/customer/%s", hostAndPort, result.ID))
			So(err, ShouldBeNil)
			So(resp.StatusCode(), ShouldEqual, 200)

			Convey("REST server should handle failed requests", func() {
				resp, _ := client.R().
					Get(fmt.Sprintf("http://%s/v1/customer/%s", hostAndPort, notExistingCustomerID))
				So(resp.StatusCode(), ShouldEqual, 404)

				Convey(fmt.Sprintf("Schedule stop signal to be sent after %s", terminateDelay), func() {
					start := time.Now()
					go func() {
						time.Sleep(terminateDelay)
						_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					}()

					Convey("Start waiting for stop signal", func() {
						s.WaitForStopSignal()

						Convey("It should wait for stop signal", func() {
							So(time.Now(), ShouldNotHappenWithin, terminateDelay, start)

							Convey("Once stop signal is received, it should call shutdown", func() {
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
	})
}
