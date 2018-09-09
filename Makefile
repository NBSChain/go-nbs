SHELL=PATH='$(PATH)' /bin/sh


PROTOC = protoc --gogofaster_out=. --proto_path=.:$(GOPATH)/src:$(dir $@) $<

# enable second expansion
.SECONDEXPANSION:


all: pbs test

test:
	go build -o nbs

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...

dir := utils/cmdKits/pb
dir2 := storage/core/pb
pbs:
	protoc -I=$(dir) --go_out=plugins=grpc:${dir} ${dir}/*.proto
	protoc -I=$(dir2) --go_out=plugins=grpc:${dir2} ${dir2}/*.proto