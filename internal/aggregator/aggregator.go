package aggregator

import (
	"fmt"
	"sync"
	"time"

	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
)

type Aggregator struct {
	Queue *queue.FlowDataQueue
	Store *store.FlowDataStore
	Wg    sync.WaitGroup
}

func NewAggregator(queue *queue.FlowDataQueue, store *store.FlowDataStore) *Aggregator {
	return &Aggregator{
		Queue: queue,
		Store: store,
	}
}

func (a *Aggregator) Worker(worker int) {
	defer a.Wg.Done()

	for flowData := range a.Queue.Channel {

		fmt.Printf("(%d): %v\n", worker, len(a.Queue.Channel))
		time.Sleep(55 * time.Millisecond)
		var key store.FlowDataTuple
		for _, item := range flowData {
			key.SrcApp = item.SrcApp
			key.DestApp = item.DestApp
			key.VpcId = item.VpcId
			key.Hour = item.Hour

			value := a.Store.LookupByTuple(key)

			value.BytesRx += item.BytesRx
			value.BytesTx += item.BytesTx

			a.Store.Save(&key, &value)
		}
	}
	fmt.Printf("end of aggregator\n")
}
