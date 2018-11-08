# nbs-node

> Full node

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)


#### Before you begin

Install [msys2] (http://www.msys2.org)

Install [grpc](https://grpc.io/docs/quickstart/go.html),[protocol-buffers](https://developers.google.com/protocol-buffers/)

go get -u github.com/golang/protobuf/protoc-gen-go

## Install
```

//shadowsocks
for windows:
	set http_proxy=http://127.0.0.1:1080
	set https_proxy=http://127.0.0.1:080
	
for linuxs:
	export http_proxy=http://127.0.0.1:1087;
	export https_proxy=http://127.0.0.1:1087;

go get -v -u github.com/NBSChain/go-nbs
cd $GOPATH/src/github.com/NBSChain/go-nbs
make
```

## Usage
```
    //start main node
    nbs
    
    //create a account(make sure the nbs service start)
    nbs account create <PASSWORD>
    
    //unlock the account
    nbs account unlock <PASSWORD>
    
    //add file in new commd window
    nbs add 1.jpg
    //get file by hash
    nbs get <HASH>
    nbs get <HASH> -o <target file name>
    //for help
    nbs --help or nbs -h
```
## Contribute
    Coming...
## License
    MIT © Protocol Labs, Inc.