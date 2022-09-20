package initiator

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
	"github.com/terracegarden/framework/config"
	datamodel "github.com/terracegarden/framework/database/model/data-model"
	"github.com/terracegarden/framework/initialize"
	opentrace "github.com/terracegarden/framework/open-trace"
	"github.com/terracegarden/framework/service"
	"github.com/uber/jaeger-client-go"
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
	tracer, closer = opentrace.Init("ScheduleService")
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

	spnCtx1 := span.Context()

	spnId := spnCtx1.(jaeger.SpanContext)
	trace_id := spnId.TraceID().String()
	span_id := spnId.SpanID().String()

	defer span.Finish()
	ctx = context.WithValue(ctx, "trace_id", trace_id)
	ctx = context.WithValue(ctx, "span_id", span_id)
	log.WithContext(ctx).Debugf("Getting schedule for %s", chipId)
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
