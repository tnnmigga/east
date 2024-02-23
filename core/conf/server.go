package conf

import "log"

var (
	serverID   uint32
	serverType string
)

func ServerID() uint32 {
	return serverID
}

func ServerType() string {
	return serverType
}

func initServerConf() {
	serverID = UInt32("server.id")
	serverType = String("server.type")
	if len(serverType) == 0 {
		log.Fatal("server.type is empty")
	}
}
