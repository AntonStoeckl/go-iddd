GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)
GO_MODULE := $(shell go mod edit -json | grep Path | head -n 1 | cut -d ":" -f 2 | cut -d '"' -f 2)
PROTO_DIR := src/customeraccounts/infrastructure/adapter/grpc
GRPC_TARGET_DIR := src/customeraccounts/infrastructure/adapter/grpc
REST_GW_TARGET_DIR := src/customeraccounts/infrastructure/adapter/rest
REST_GW_OUT_FILE := customer.pb.gw.go

generate_proto:
	@protoc \
		-I $(GRPC_TARGET_DIR) \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:$(GRPC_TARGET_DIR) \
		--grpc-gateway_out=logtostderr=true,import_path=customerrest:$(REST_GW_TARGET_DIR) \
		--swagger_out=logtostderr=true:$(REST_GW_TARGET_DIR) \
		$(PROTO_DIR)/customer.proto

	@# Not possible to split grpc and rest otherwise: https://github.com/grpc-ecosystem/grpc-gateway/issues/353
	@sed -i '/package customerrest/ a \\nimport customergrpc "$(GO_MODULE)/$(GRPC_TARGET_DIR)"' $(REST_GW_TARGET_DIR)/$(REST_GW_OUT_FILE)
	@sed -i 's/client CustomerClient/client customergrpc.CustomerClient/' $(REST_GW_TARGET_DIR)/$(REST_GW_OUT_FILE)
	@sed -i 's/server CustomerServer/server customergrpc.CustomerServer/' $(REST_GW_TARGET_DIR)/$(REST_GW_OUT_FILE)
	@sed -i 's/NewCustomerClient/customergrpc.NewCustomerClient/' $(REST_GW_TARGET_DIR)/$(REST_GW_OUT_FILE)
	@sed -i -E 's/var protoReq (.+)/var protoReq customergrpc.\1/' $(REST_GW_TARGET_DIR)/$(REST_GW_OUT_FILE)

lint:
	golangci-lint run --build-tags test ./...

# https://github.com/golangci/golangci-lint
install-golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.24.0


# https://github.com/psampaz/go-mod-outdated
outdated-list:
	go list -u -m -json all | go-mod-outdated -update -direct