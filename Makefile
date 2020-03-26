GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I service/customer/infrastructure/primary/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:service/customer/infrastructure/primary/grpc \
		--grpc-gateway_out=logtostderr=true:service/customer/infrastructure/primary/grpc \
		--swagger_out=logtostderr=true:service/customer/infrastructure/primary/grpc \
		service/customer/infrastructure/primary/grpc/customer.proto

generate_mocked_EventStore:
	@mockery \
		-name EventStore \
		-dir service/lib/es \
		-outpkg mocked \
		-output service/lib/eventstore/mocked \
		-note "+build test"

generate_mocked_ForStoringCustomerEvents:
	@mockery \
		-name ForStoringCustomerEvents \
		-dir service/customer/application/command \
		-outpkg mocked \
		-output service/customer/infrastructure/secondary/mocked \
		-note "+build test"

generate_mocked_ForAssertingUniqueEmailAddresses:
	@mockery \
		-name ForAssertingUniqueEmailAddresses \
		-dir service/customer/infrastructure/secondary/eventstore \
		-outpkg mocked \
		-output service/customer/infrastructure/secondary/mocked \
		-note "+build test"

generate_all_mocks: \
	generate_mocked_EventStore \
	generate_mocked_ForStoringCustomerEvents \
	generate_mocked_ForAssertingUniqueEmailAddresses

lint:
	golangci-lint run --build-tags test ./...

# https://github.com/golangci/golangci-lint
install-golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.23.8


# https://github.com/psampaz/go-mod-outdated
outdated-list:
	go list -u -m -json all | go-mod-outdated -update -direct