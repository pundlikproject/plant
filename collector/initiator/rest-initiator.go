package initiator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	plantservice "github.com/terracegarden/collector/pkg/plant"
	"github.com/terracegarden/framework/config"
	datamodel "github.com/terracegarden/framework/database/model/data-model"
	"github.com/terracegarden/framework/initialize"
	internal "github.com/terracegarden/framework/internal-bus"
	"github.com/terracegarden/framework/service"
	jaeger "github.com/uber/jaeger-client-go"
	configJ "github.com/uber/jaeger-client-go/config"
)

var tracer opentracing.Tracer
var closer io.Closer

func InitFramework(ctx context.Context) {
	//Init logger
	tracer, closer = Init("plant-chip")

	//defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	cfg := config.FrameworkConfig{
		DbConf: &config.DbConfig{Url: "localhost:5432", Database: "postgres", UserName: "postgres", Password: "passwd@123", PoolSize: 100},
	}
	err := initialize.Init(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func StartRestServer(ctx context.Context) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rest/save_user", saveUser).Methods("POST")
	router.HandleFunc("/rest/save_chip", saveChip).Methods("POST")
	router.HandleFunc("/rest/save_plant", savePlant).Methods("POST")

	router.HandleFunc("/rest/get_chips/{ownerId}", GetChips).Methods("GET")
	router.HandleFunc("/rest/get_plants/{chipId}", GetPlants).Methods("GET")
	router.HandleFunc("/rest/set_schedule/{chipId}", SetSchedule).Methods("POST")

	log.Println("Starting server.")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func getBody(ctx context.Context, obj interface{}, r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

func saveUser(w http.ResponseWriter, r *http.Request) {
	owt := &datamodel.Owner{}
	err := getBody(r.Context(), owt, r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = service.RegisterUser(r.Context(), owt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func saveChip(w http.ResponseWriter, r *http.Request) {
	chp := &datamodel.Chip{}
	err := getBody(r.Context(), chp, r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = service.SaveChip(r.Context(), chp)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func savePlant(w http.ResponseWriter, r *http.Request) {
	plnt := &datamodel.Plant{}
	err := getBody(r.Context(), plnt, r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = service.SavePlant(r.Context(), plnt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func GetChips(w http.ResponseWriter, r *http.Request) {
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	vars := mux.Vars(r)
	ownerId, ok := vars["ownerId"]
	if !ok {
		return
	}

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "GetChips", ext.RPCServerOption(spanCtx))
	span.LogFields(
		traceLog.String("owner-id", ownerId),
	)
	defer span.Finish()

	chips, err := service.GetChips(ctx, ownerId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error code %v", err), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(chips)
}

func GetPlants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chipId, ok := vars["chipId"]
	if !ok {
		return
	}

	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "GetPlants", ext.RPCServerOption(spanCtx))
	span.LogFields(
		traceLog.String("chip-id", chipId),
	)
	defer span.Finish()

	plants, err := plantservice.GetPlants(ctx, chipId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error code %v", err), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(plants)

}

func SetSchedule(w http.ResponseWriter, r *http.Request) {
	schedules := make([]*datamodel.Schedule, 0)
	err := getBody(r.Context(), schedules, r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	schResp, err := internal.CallInternalService(r.Context(), internal.InternalServiceSchedule, []interface{}{schedules[0].ChipId}, internal.MethodPost, schedules)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	rsp, _ := json.Marshal(schResp)
	w.Write(rsp)
}

func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &configJ.Configuration{
		ServiceName: service,
		Sampler: &configJ.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &configJ.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(configJ.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
