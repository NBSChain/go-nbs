# nbs-node

> Full node

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)


###Before you begin
Install [grpc](https://grpc.io/docs/quickstart/go.html)    

## Install
```


//shadowsocks
for windows:
	set http_proxy=http://127.0.0.1:1080
	set https_proxy=http://127.0.0.1:080
	
for linuxs:
	export http_proxy=http://127.0.0.1:1087;
	export https_proxy=http://127.0.0.1:1087;

go get github.com/NBSChain/go-nbs
cd $GOPATH/src/github.com/NBSChain/go-nbs
make
```

## Usage
```
    //start main node
    nbs
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