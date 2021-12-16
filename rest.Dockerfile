FROM golang:1.17.3-alpine as build-env

WORKDIR $GOPATH/src/github.com/AntonStoeckl/go-iddd
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goapp ./src/service/rest/cmd

FROM alpine:3
RUN mkdir /app
# Create user and set ownership and permissions as required
RUN adduser -D myuser && chown -R myuser /app
WORKDIR /app
USER myuser
COPY --from=build-env /go/src/github.com/AntonStoeckl/go-iddd/goapp .

EXPOSE 8085
ENTRYPOINT ["./goapp"]