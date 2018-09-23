SHELL=PATH='$(PATH)' /bin/sh


PROTOC = protoc --gogofaster_out=. --proto_path=.:$(GOPATH)/src:$(dir $@) $<

# enable second expansion
.SECONDEXPANSION:


all: pbs build

build:
	go build -o nbs

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...

dir := utils/cmdKits/pb
dir2 := storage/application/pb
dir3 := storage/merkledag/pb

pbs:
	protoc -I=$(dir) --go_out=plugins=grpc:${dir} ${dir}/*.proto
	protoc -I=$(dir2) --go_out=plugins=grpc:${dir2} ${dir2}/*.proto
	protoc -I=$(dir3) -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf \
	--go_out=plugins=grpc:${dir3} ${dir3}/*.proto

clean:
	rm -rf nbs

test:
	go test -v ./storage/application/rpcService/