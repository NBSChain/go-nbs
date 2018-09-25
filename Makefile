SHELL=PATH='$(PATH)' /bin/sh

PLATFORM := $(shell uname -o)


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
	go build -o $(EXENAME)

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...

dir := utils/cmdKits/pb
dir2 := storage/application/pb
dir3 := storage/merkledag/pb

pbs:
	protoc -I=$(dir) --go_out=plugins=grpc:${dir} ${dir}/*.proto
	protoc -I=$(dir2) --go_out=plugins=grpc:${dir2} ${dir2}/*.proto
	protoc -I=$(dir3) -I=$(INCLUDE)/src -I=$(INCLUDE)/src/github.com/gogo/protobuf/protobuf --go_out=plugins=grpc:${dir3} ${dir3}/*.proto

clean:
	rm -rf nbs

test:
	go test -v ./storage/application/rpcService/
