# echo product.proto
grpc-proto:
	protoc -I. --go_out=plugins=grpc:. proto/product/product.proto

# grpc protoset
grpc-protoset:
	protoc --proto_path=./ \
	--descriptor_set_out=product.protoset \
	--include_imports \
	proto/product/product.proto

#grpc-gateway
grpc-gateway:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/product/product.yml:. \
	proto/product/product.proto	


# Generate swagger definitions
grpc-swagger:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--swagger_out=logtostderr=true,grpc_api_configuration=proto/product/product.yml:. \
	proto/product/product.proto
   
# Command to create the module
mod:
	go mod init product-service.com/product-service
# echo Dockerfile
build:
	docker build -f Dockerfile -t product-service:1 .
# run
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 product-service:1
