package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/mmcken3/snap-plugin-publisher-udpsocket/udpsocket"
)

func main() {
	meta := UDPSocket.Meta()
	plugin.Start(meta, UDPSocket.NewUdpPublisher(), os.Args[1])
}
