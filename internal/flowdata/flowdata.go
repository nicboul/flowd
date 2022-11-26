package flowdata

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/nicboul/flowdata/internal/flowdataread"
	"github.com/nicboul/flowdata/internal/flowdatawrite"
	"github.com/nicboul/flowdata/internal/store"
)

type FlowDataParams struct {
	Timeout time.Duration
	Store   *store.FlowDataStore
}

func NewFlowDataServer(p FlowDataParams) *mux.Router {

	p.Store = store.NewFlowDataStore()

	flowDataWriteHandler := &flowdatawrite.FlowDataWrite{
		Store: p.Store,
	}

	flowDataReadHandler := &flowdataread.FlowDataRead{
		Store: p.Store,
	}

	muxRouter := mux.NewRouter()
	muxRouter.Methods("POST").PathPrefix("/flows").Handler(flowDataWriteHandler)
	muxRouter.Methods("GET").PathPrefix("/flows").Handler(flowDataReadHandler)

	return muxRouter
}
