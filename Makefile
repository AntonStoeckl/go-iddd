GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I service/customer/infrastructure/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:service/customer/infrastructure/grpc \
		--grpc-gateway_out=logtostderr=true:service/customer/infrastructure/grpc \
		--swagger_out=logtostderr=true:service/customer/infrastructure/grpc \
		service/customer/infrastructure/grpc/customer.proto