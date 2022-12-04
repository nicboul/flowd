package flowdata

import (
	"context"
	"errors"
	"net/http"
	"sync"
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

func scaleAggregator(s *FlowDataServer, max int) {
	w := 0
	for {
		if w < max && len(s.Params.Queue.Channel) > 20 {
			s.aggregator.Wg.Add(1)
			w++
			go s.aggregator.Worker(w)

		}
		time.Sleep(1 * time.Second)
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (s *FlowDataServer) Serve(ctx context.Context) error {

	httpSrv := &http.Server{
		Addr:    s.Params.Listen,
		Handler: s.MuxRouter,
	}

	errGroup, errctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		<-errctx.Done()
		return httpSrv.Close()
	})

	errGroup.Go(func() error {
		err := httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	go scaleAggregator(s, 100)

	/* Wait for the http server to stop running */
	err := errGroup.Wait()

	/* Close the Queue (the channel being needs to be closed */
	s.Params.Queue.Close()
	/* Wait for the workers to finish their job */
	waitTimeout(&s.aggregator.Wg, 1*time.Minute)

	return err
}
