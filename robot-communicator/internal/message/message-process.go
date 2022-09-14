package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mochi-co/mqtt/server/events"
	datamodel "github.com/terracegarden/framework/database/model/data-model"
	"github.com/terracegarden/framework/service"
)

func ProcessMessage(cl events.Client, pk events.Packet) (events.Packet, error) {
	fmt.Printf("OnMessage %s: %+v\n", cl.ID, pk.Payload)
	ctx := context.Background()
	data, err := unmarshalMessage(ctx, pk.TopicName, pk.Payload)
	if err != nil {
		return pk, err
	}

	plantData, ok := data.(*datamodel.Plant)
	if !ok {
		return pk, errors.New("unwanted error type")
	}
	err = service.SavePlant(ctx, plantData)
	return pk, err
}

func unmarshalMessage(ctx context.Context, topic string, data []byte) (interface{}, error) {
	plnt := &datamodel.Plant{}
	err := json.Unmarshal(data, plnt)
	if err != nil {
		return nil, err
	}
	return plnt, nil
}
