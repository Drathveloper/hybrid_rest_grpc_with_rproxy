GOPATH=$(HOME)/go
PROTOC_VERSION=23.4
PROTOC_GEN_GO_VERSION=1.31
ARCH_VERSION=$(shell uname -m)
OS_VERSION=$(shell uname)
PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(OS_VERSION)-$(ARCH_VERSION).zip"
PROTOC_PATH=$(CURDIR)/tools

generate-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $(CURDIR)/certs/nginx-selfsigned.key -out $(CURDIR)/certs/nginx-selfsigned.crt -batch -config $(CURDIR)/cert.conf

generate-grpc:
	echo $(PATH)
	mkdir -p $(CURDIR)/tools
	curl -sSfL $(PROTOC_URL) -o $(PROTOC_PATH)/protoc
	chmod a+x $(PROTOC_PATH)/protoc
	cd $(PROTOC_PATH)
	export PATH="$(PATH):$(GOPATH)/bin"; protoc --go_out=. --go-grpc_out=. ./resources/*.proto

compile: generate-grpc
	mkdir -p $(CURDIR)/bin
	go build -o $(CURDIR)/bin/poc_hybrid_grpc_rest