package flowdata

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/nicboul/flowdata/internal/aggregator"
	"github.com/nicboul/flowdata/internal/flowdataread"
	"github.com/nicboul/flowdata/internal/flowdatawrite"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
)

type FlowDataParams struct {
	Timeout time.Duration
	Store   *store.FlowDataStore
	Queue   *queue.FlowDataQueue
}

type FlowDataServer struct {
	aggregator *aggregator.Aggregator
	writer     *flowdatawrite.FlowDataWrite
	reader     *flowdataread.FlowDataRead
	MuxRouter  *mux.Router
}

func NewFlowDataServer(p FlowDataParams) *FlowDataServer {

	server := FlowDataServer{}

	server.aggregator = aggregator.NewAggregator(p.Queue, p.Store)
	go server.aggregator.Worker()

	server.writer = flowdatawrite.NewFlowDataWrite(p.Queue)
	server.reader = flowdataread.NewFlowDataRead(p.Store)

	server.MuxRouter = mux.NewRouter()
	server.MuxRouter.Methods("POST").PathPrefix("/flows").Handler(server.writer)
	server.MuxRouter.Methods("GET").PathPrefix("/flows").Handler(server.reader)

	return &server
}
