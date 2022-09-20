package initiator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
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

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	if entry.Context != nil {
		entry.Data["trace_id"] = entry.Context.Value("trace_id")
		entry.Data["span_id"] = entry.Context.Value("span_id")
	}
	return l.formatter.Format(entry)
}

func InitFramework(ctx context.Context) {
	log.SetLevel(log.DebugLevel)

	//Init logger
	tracer, closer = Init("PlantService")

	//defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})

	cfg := config.FrameworkConfig{
		DbConf: &config.DbConfig{Url: "postgres:5432", Database: "postgres", UserName: "postgres", Password: "passwd@123", PoolSize: 100},
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

	log.Debug("Starting server.")
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
	spnCtx1 := span.Context()

	spnId := spnCtx1.(jaeger.SpanContext)
	trace_id := spnId.TraceID().String()
	span_id := spnId.SpanID().String()
	defer span.Finish()
	ctx = context.WithValue(ctx, "trace_id", trace_id)
	ctx = context.WithValue(ctx, "span_id", span_id)

	log.WithContext(ctx).Debugln("Function")
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
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Get-Plants", ext.RPCServerOption(spanCtx))
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
	cfg, _ := configJ.FromEnv()
	cfg.ServiceName = service
	cfg.Reporter.LogSpans = true
	log.Infof("Jeager config %+v", cfg)
	tracer, closer, err := cfg.NewTracer(configJ.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
