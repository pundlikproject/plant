package main

import (
	"context"
	"log"

	"github.com/terracegarden/collector/initiator"
)

func main() {
	ctx := context.Background()

	//Initialize framework
	initiator.InitFramework(ctx)

	//log.Print("Starting MQTT server")
	//Start mqtt server
	//mqttserv.StartMqttServer(ctx, ":1883")

	log.Print("Starting REST server")
	//Starting rest server
	initiator.StartRestServer(ctx)
}
