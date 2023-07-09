# hybrid_rest_grpc_with_rproxy

Golang project to test the possibility of have custom reverse proxy in front of a service that exposes both gRPC and REST endpoints.

# Setup:
sudo apt install protoc

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# Generate protocol buffers code:
make generate-grpc

OR alternatively if make fails, you can do it yourself by hand with this command:

protoc --go_out=. --go-grpc_out=. ./resources/*.proto

# Generate SSL certs:
make generate-certs

# Run:
docker-compose up
The server is not directly accesible, and nginx is accesible through port 443
