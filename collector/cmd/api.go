package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/terracegarden/collector/initiator"
	mqttserv "github.com/terracegarden/collector/internal/mqtt-server"
)

func main() {
	ctx := context.Background()

	//Initialize framework
	initiator.InitFramework(ctx)

	log.Info("Starting MQTT server")
	//Start mqtt server
	mqttserv.StartMqttServer(ctx, ":1883")

	log.Info("Starting REST server")
	//Starting rest server
	initiator.StartRestServer(ctx)
}
