package initiator

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"github.com/terracegarden/framework/config"
	datamodel "github.com/terracegarden/framework/database/model/data-model"
	"github.com/terracegarden/framework/initialize"
	opentrace "github.com/terracegarden/framework/open-trace"
	"github.com/terracegarden/framework/service"
)

var tracer opentracing.Tracer
var closer io.Closer

func InitFramework(ctx context.Context) {
	tracer, closer = opentrace.Init("ScheduleService")

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
	router.HandleFunc("/rest/schedule/{chipId}", getSchedule).Methods("GET")
	router.HandleFunc("/rest/schedule/{chipId}", saveSchedule).Methods("POST")

	log.Println("Starting server....")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBody(ctx context.Context, obj interface{}, r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	vars := mux.Vars(r)
	chipId, ok := vars["chipId"]
	if !ok {
		w.Write([]byte("chipId is mandatory"))
		return
	}

	span, ctx := opentracing.StartSpanFromContext(r.Context(), "GetChips", ext.RPCServerOption(spanCtx))
	span.LogFields(
		traceLog.String("chipId", chipId),
	)
	defer span.Finish()

	schedules, err := service.GetSchedule(ctx, chipId)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	schByt, _ := json.Marshal(schedules)
	w.Write(schByt)
}

func saveSchedule(w http.ResponseWriter, r *http.Request) {
	chp := &datamodel.Schedule{}
	err := getBody(r.Context(), chp, r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = service.SaveSchedule(r.Context(), chp)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
