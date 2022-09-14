package mqttserver

import (
	"context"
	"fmt"
	"log"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/terracegarden/collector/internal/message"
)

func StartMqttServer(ctx context.Context, port string) error {
	// Create the new MQTT Server.
	server := mqtt.NewServer(nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883")

	// Add the listener to the server with default options (nil).
	err := server.AddListener(tcp, nil)
	if err != nil {
		log.Fatal(err)
	}

	server.Events.OnConnect = func(cl events.Client, pk events.Packet) {
		fmt.Printf("<< OnConnect client connected %s: %+v\n", cl.ID, pk)
	}

	server.Events.OnMessage = message.ProcessMessage

	// Start the broker. Serve() is blocking - see examples folder
	// for usage ideas.
	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
