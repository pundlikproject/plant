package plant

import (
	"context"
	"encoding/json"

	datamodel "github.com/terracegarden/framework/database/model/data-model"
	internal "github.com/terracegarden/framework/internal-bus"
	"github.com/terracegarden/framework/service"
)

func GetPlants(ctx context.Context, chipId string) ([]*datamodel.Plant, error) {
	plants, err := service.GetPlants(ctx, chipId)
	if err != nil {
		return nil, err
	}
	_ = plants

	scheduleByt, err := internal.CallInternalService(ctx, internal.InternalServiceSchedule, []interface{}{chipId}, internal.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	schedules := make([]*datamodel.Schedule, 0)
	json.Unmarshal(scheduleByt, &schedules)
	scheduleMap := map[string]*datamodel.Schedule{}

	for _, sch := range schedules {
		scheduleMap[sch.PlantId] = sch
	}

	for _, Plnt := range plants {
		if _, present := scheduleMap[Plnt.PlantId]; present {
			Plnt.Schedule = scheduleMap[Plnt.PlantId].Schedule
		}
	}

	return plants, nil
}
