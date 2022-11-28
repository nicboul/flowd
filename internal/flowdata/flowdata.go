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

func NewFlowDataServer(p FlowDataParams) *mux.Router {

	aggregator := aggregator.NewAggregator(p.Queue, p.Store)
	go aggregator.Aggregator()

	flowDataWriteHandler := &flowdatawrite.FlowDataWrite{
		Queue: p.Queue,
	}

	flowDataReadHandler := &flowdataread.FlowDataRead{
		Store: p.Store,
	}

	muxRouter := mux.NewRouter()
	muxRouter.Methods("POST").PathPrefix("/flows").Handler(flowDataWriteHandler)
	muxRouter.Methods("GET").PathPrefix("/flows").Handler(flowDataReadHandler)

	return muxRouter
}
