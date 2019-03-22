SHELL=PATH='$(PATH)' /bin/sh

ifeq ($(OS),"")
    PLATFORM := $(shell uname -o)
else
    PLATFORM := $(shell uname -s)
endif

EXTEND := .exe
ifeq ($(PLATFORM), Msys)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else
	INCLUDE := $(GOPATH)
	INCLUDESRC := ${shell echo "$(GOPATH)"|sed -e 's/:/\/src:/g'}
	INCLUDEBIN := ${shell echo "$(GOPATH)"|sed -e 's/:/\/bin:/g'}
	INCLUDESRC := $(INCLUDESRC)/src
	INCLUDEBIN := $(INCLUDEBIN)/bin
    FIRSTBIN:=${shell echo "$(INCLUDEBIN)"|cut -d':' -f1}
	EXTEND	:=
endif

include CygwinMF.mk

EXENAME := nbs$(EXTEND)

# enable second expansion
.SECONDEXPANSION:

all: pbs build

build:
	go build -race -o $(EXENAME)
	mv $(EXENAME) $(FIRSTBIN)

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
	protoc -I=${gossip} -I=$(INCLUDESRC)   --go_out=plugins=grpc:${gossip} 		${gossip}/*.proto

clean:
	rm -rf $(FIRSTBIN)/nbs

test:
	go test -v ./storage/application/rpcService/
	go test -v ./utils/crypto/

