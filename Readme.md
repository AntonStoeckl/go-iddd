# Work in progress!

Once finished, this should server as a showcase how I implement a DDD style project in Go(lang),
with a Ports & Adapter architecture and EventSourcing.

I plan to use it for a series of blog posts about this topic.

### Setup for local development

#### Start Docker container(s)

Run `docker-compose up -d` in the project root.

#### Environment configuration

##### To be able to start the service

Create local.env file in the project root (.env files is gitignored there) with following contents and replace
$PathToProjectRoot$ with the path to the root of the go-idd project sources.

```
POSTGRES_DSN=postgresql://goiddd:password123@localhost:15432/goiddd_local?sslmode=disable
POSTGRES_MIGRATIONS_PATH=$PathToProjectRoot$/go-iddd/service/lib/eventstore/postgres/database/migrations
```

##### To be able to run the tests

Create test.env file in the project root (.env files is gitignored there) with following contents and replace
$PathToProjectRoot$ with the path to the root of the go-idd project sources.

```
POSTGRES_DSN=postgresql://goiddd:password123@localhost:15432/goiddd_test?sslmode=disable
POSTGRES_MIGRATIONS_PATH=$PathToProjectRoot$/go-iddd/service/lib/eventstore/postgres/database/migrations
```

##### To run HTTP requests with GoLand's (IntelliJ) new built-in HTTP client

Create customer.http file in the project root (.http files are gitignored there) with following contents.

```
### Register a Customer
POST http://localhost:8085/v1/customer
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "emailAddress": "john@doe.com",
  "familyName": "Doe",
  "givenName": "John"
}

> {% client.global.set("id", response.body.id); %}

### Confirm a Customer's email address
PUT http://localhost:8085/v1/customer/{{id}}/emailaddress/confirm
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "confirmationHash": "0acf14bbeaf0b9c6ef8e39d7f9254336"
}

### Change a Customer's email address
PUT http://localhost:8085/v1/customer/{{id}}/emailaddress
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "emailAddress": "john+changed@doe.com"
}

### Delete a Customer
DELETE http://localhost:8085/v1/customer/{{id}}
Accept: application/json
Cache-Control: no-cache
Content-Type: application/json

### Get the Swagger documentation
GET http://localhost:8085/v1/customer/swagger.json

###
```

**Attention**

The *ConfirmEmailAddress* request does not work without changes - the *confirmationHash* needs to be adapted.
You can find it in the *CustomerRegistered* event in the eventstore DB table.
For security reasons the hash not returned in the response of the *Register* request (it **must** only be sent to the Customer via email ;-)

#### Start the service (gRPC and REST)

##### Via Terminal

1) Source the local.env file in your terminal, e.g. `source dev/local.env` or set the env vars in a different way
2) In the project root run `go run service/cmd/grpc/main.go`

##### Via GoLand

1) Create a build configuration for `service/cmd/grpc/main.go`
2) I suggest to use the great GoLand plugin [EnvFile](https://plugins.jetbrains.com/plugin/7861-envfile)
and add the local.env file in the build configuration
