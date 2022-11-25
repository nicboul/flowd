package aggregator

import (
	"github.com/nicboul/flowdata/internal/flowdata"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
)

func Aggregate() {
	var flowData []flowdata.FlowData
	flowData = queue.Consume()

	var key store.FlowDataKey

	for _, item := range flowData {
		key.SrcApp = item.SrcApp
		key.DestApp = item.DestApp
		key.VpcId = item.VpcId
		key.Hour = item.Hour

		value := store.LookupValue(key)

		value.BytesRx += item.BytesRx
		value.BytesTx += item.BytesTx

		store.Save(key, value)
	}
}
