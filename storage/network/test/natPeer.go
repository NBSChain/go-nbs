package main

import (
	"flag"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network"
	"github.com/NBSChain/go-nbs/thirdParty/idService"
)

func main() {

	targetID := flag.String("d", "", "server ip")
	localPort := flag.Int("p", 9900, "local server listening port")
	flag.Parse()

	netInstance := network.GetInstance()

	if *targetID != "" {
		host := netInstance.NewHost()
		host.SetStreamHandler("/test/1.0.0", func(s network.Stream) {

		})
		host.NewStream(*targetID, "/test/1.0.0")
	} else {

		l := fmt.Sprintf("%d", localPort)

		id, _ := idService.GetInstance().GenerateId("")

		netInstance.NewHost(
			netInstance.Identity(id),
			netInstance.ListenAddrString(l),
		)

		<-make(chan struct{})
	}
}
