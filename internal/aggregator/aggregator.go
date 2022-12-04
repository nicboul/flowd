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

		//fmt.Printf("(%d): %v\n", worker, len(a.Queue.Channel))
		time.Sleep(300 * time.Millisecond)
		var key store.FlowDataTuple
		for _, item := range flowData {
			key.SrcApp = item.SrcApp
			key.DestApp = item.DestApp
			key.VpcId = item.VpcId
			key.Hour = item.Hour

			/* Lock the whole transaction */
			a.Store.Lock.Lock()
			value := a.Store.LookupByTupleWithLock(key)

			value.BytesRx += item.BytesRx
			value.BytesTx += item.BytesTx

			a.Store.SaveWithLock(&key, &value)
			a.Store.Lock.Unlock()
		}
	}
	fmt.Printf("end of aggregator (%d)\n", worker)
}
