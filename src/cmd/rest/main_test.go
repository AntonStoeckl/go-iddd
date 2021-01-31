package main

import (
	"context"
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/service"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func TestStartRestServer(t *testing.T) {
	var restServer *http.Server
	var grpcClientConn *grpc.ClientConn

	//logger := shared.NewStandardLogger()
	logger := shared.NewNilLogger()
	config := service.MustBuildConfigFromEnv(logger)
	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)

	exitWasCalled := false
	exit := func() {
		exitWasCalled = true
	}
	myShutdown := func() {
		shutdown(logger, cancelCtx, grpcClientConn, restServer, exit)
	}

	restServer, grpcClientConn = buildRestServer(config, logger, ctx, myShutdown)

	terminateDelay := time.Millisecond * 100

	Convey("Start the REST server as a goroutine", t, func() {
		go startRestServer(config, logger, restServer, myShutdown)

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
						waitForStopSignal(logger, myShutdown)

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
