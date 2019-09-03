GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate:
	@protoc \
		-I api/grpc/customer \
		-I /usr/local/include/google/protobuf/ \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:api/grpc/customer \
		api/grpc/customer/customer.proto