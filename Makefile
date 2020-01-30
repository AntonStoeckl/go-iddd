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
		-dir service/customer/application \
		-outpkg mocked \
		-output service/customer/infrastructure/secondary/forstoringcustomerevents/mocked \
		-note "+build test"

generate_all_mocks: \
	generate_mocked_EventStore \
	generate_mocked_ForStoringCustomerEvents
