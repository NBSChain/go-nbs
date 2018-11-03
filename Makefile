SHELL=PATH='$(PATH)' /bin/sh

PLATFORM := $(shell uname -a)


EXTEND := .exe
ifeq ($(PLATFORM), Msys)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else
	INCLUDE := $(GOPATH)
	EXTEND	:=
endif

EXENAME := nbs$(EXTEND)

# enable second expansion
.SECONDEXPANSION:

	echo $(PLATFORM)

all: pbs build

build:
	go build -race -o $(EXENAME)
	mv $(EXENAME) $(GOPATH)/bin/

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...

dir := console/pb
dir2 := storage/application/pb
dir3 := storage/merkledag/pb
dir4 := storage/bitswap/pb
dir5 := storage/network/pb
dir6 := thirdParty/account/pb

pbs:
	protoc -I=$(dir) --go_out=plugins=grpc:${dir} ${dir}/*.proto
	protoc -I=$(dir2) --go_out=plugins=grpc:${dir2} ${dir2}/*.proto
	protoc -I=$(dir4) --go_out=plugins=grpc:${dir4} ${dir4}/*.proto
	protoc -I=$(dir5) --go_out=plugins=grpc:${dir5} ${dir5}/*.proto
	protoc -I=$(dir6) --go_out=plugins=grpc:${dir6} ${dir6}/*.proto
	protoc -I=$(dir3) -I=$(INCLUDE)/src -I=$(INCLUDE)/src/github.com/gogo/protobuf/protobuf --go_out=plugins=grpc:${dir3} ${dir3}/*.proto

clean:
	rm -rf nbs

test:
	go test -v ./storage/application/rpcService/
	go test -v ./utils/crypto/

