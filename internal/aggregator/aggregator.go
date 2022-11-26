package aggregator

import (
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
)

func Aggregate(s *store.FlowDataStore) {
	var flowData []queue.FlowData
	flowData = queue.Consume()

	var key store.FlowDataKey

	for _, item := range flowData {
		key.SrcApp = item.SrcApp
		key.DestApp = item.DestApp
		key.VpcId = item.VpcId
		key.Hour = item.Hour

		value := s.LookupValue(key)

		value.BytesRx += item.BytesRx
		value.BytesTx += item.BytesTx

		s.Save(key, value)
	}
}
