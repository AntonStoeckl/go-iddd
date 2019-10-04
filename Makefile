GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I customer/api/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:customer/api/grpc \
		--grpc-gateway_out=logtostderr=true:customer/api/grpc \
		--swagger_out=logtostderr=true:customer/api/grpc \
		customer/api/grpc/customer.proto