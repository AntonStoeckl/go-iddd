generate:
	@protoc \
		-I api/grpc/customer \
		-I /usr/local/include/google/protobuf/ \
		--go_out=plugins=grpc:api/grpc/customer \
		api/grpc/customer/customer.proto