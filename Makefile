GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I customer/infrastructure/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:customer/infrastructure/grpc \
		--grpc-gateway_out=logtostderr=true:customer/infrastructure/grpc \
		--swagger_out=logtostderr=true:customer/infrastructure/grpc \
		customer/infrastructure/grpc/customer.proto