GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I service/customeraccounts/infrastructure/adapter/primary/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:service/customeraccounts/infrastructure/adapter/primary/grpc \
		--grpc-gateway_out=logtostderr=true:service/customeraccounts/infrastructure/adapter/primary/grpc \
		--swagger_out=logtostderr=true:service/customeraccounts/infrastructure/adapter/primary/grpc \
		service/customeraccounts/infrastructure/adapter/primary/grpc/customer.proto

lint:
	golangci-lint run --build-tags test ./...

# https://github.com/golangci/golangci-lint
install-golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.23.8


# https://github.com/psampaz/go-mod-outdated
outdated-list:
	go list -u -m -json all | go-mod-outdated -update -direct