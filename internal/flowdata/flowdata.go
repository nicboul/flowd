package flowdata

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicboul/flowdata/internal/aggregator"
	"github.com/nicboul/flowdata/internal/flowdataread"
	"github.com/nicboul/flowdata/internal/flowdatawrite"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
	"golang.org/x/sync/errgroup"
)

type FlowDataParams struct {
	Timeout time.Duration
	Store   *store.FlowDataStore
	Queue   *queue.FlowDataQueue
	Listen  string
}

type FlowDataServer struct {
	aggregator *aggregator.Aggregator
	writer     *flowdatawrite.FlowDataWrite
	reader     *flowdataread.FlowDataRead
	MuxRouter  *mux.Router
	Params     *FlowDataParams
}

func NewFlowDataServer(p FlowDataParams) *FlowDataServer {

	server := FlowDataServer{}
	server.Params = &p

	server.aggregator = aggregator.NewAggregator(p.Queue, p.Store)

	server.writer = flowdatawrite.NewFlowDataWrite(p.Queue)
	server.reader = flowdataread.NewFlowDataRead(p.Store)

	server.MuxRouter = mux.NewRouter()
	server.MuxRouter.Methods("POST").PathPrefix("/flows").Handler(server.writer)
	server.MuxRouter.Methods("GET").PathPrefix("/flows").Handler(server.reader)

	return &server
}

func (s *FlowDataServer) Serve(ctx context.Context) error {

	httpSrv := &http.Server{
		Addr:    s.Params.Listen,
		Handler: s.MuxRouter,
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		<-ctx.Done()
		return httpSrv.Close()
	})

	errGroup.Go(func() error {
		err := httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	go s.aggregator.Worker()

	return errGroup.Wait()

}
