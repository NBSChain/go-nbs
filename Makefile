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
	go build -race -o $(EXENAME)
	mv $(EXENAME) $(INCLUDE)/bin/

deps:
	go get -u -d -v github.com/libp2p/go-libp2p/...

console := console/pb
application := storage/application/pb
ipld := storage/merkledag/pb
bitswap := storage/bitswap/pb
network := storage/network/pb
account := thirdParty/account/pb
gossip := thirdParty/gossip/pb

pbs:
	protoc -I=$(console)  		--go_out=plugins=grpc:${console} 		${console}/*.proto
	protoc -I=$(application) 	--go_out=plugins=grpc:${application} 	${application}/*.proto
	protoc -I=$(ipld) 			--go_out=plugins=grpc:${ipld} 			${ipld}/*.proto
	protoc -I=$(bitswap) 		--go_out=plugins=grpc:${bitswap} 		${bitswap}/*.proto
	protoc -I=$(network) 		--go_out=plugins=grpc:${network} 		${network}/*.proto
	protoc -I=$(account) 		--go_out=plugins=grpc:${account} 		${account}/*.proto
	protoc -I=$(gossip) 		--go_out=plugins=grpc:${gossip} 		${gossip}/*.proto

clean:
	rm -rf nbs

test:
	go test -v ./storage/application/rpcService/
	go test -v ./utils/crypto/

