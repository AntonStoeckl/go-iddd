generate:
	@protoc \
		-I api/rpc/grpc/customer \
		-I /usr/local/include/google/protobuf/ \
		--go_out=plugins=grpc:api/rpc/grpc/customer \
		api/rpc/grpc/customer/customer.proto