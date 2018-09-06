SHELL=PATH='$(PATH)' /bin/sh


PROTOC = protoc --gogofaster_out=. --proto_path=.:$(GOPATH)/src:$(dir $@) $<

# enable second expansion
.SECONDEXPANSION:


test:
	go build -o nbs

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...
