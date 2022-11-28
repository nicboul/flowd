package flowdataread

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
	log "github.com/sirupsen/logrus"
)

type FlowDataRead struct {
	Store *store.FlowDataStore
}

func (f *FlowDataRead) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hourStr := r.URL.Query().Get("hour")
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	flowDataResponse := []queue.FlowData{}

	flows := f.Store.LookupByHour(hour)
	for key, value := range flows {
		var flowData queue.FlowData
		flowData.SrcApp = key.SrcApp
		flowData.DestApp = key.DestApp
		flowData.VpcId = key.VpcId
		flowData.Hour = key.Hour

		flowData.BytesRx = value.BytesRx
		flowData.BytesTx = value.BytesTx

		flowDataResponse = append(flowDataResponse, flowData)
	}

	jsonStr, err := json.Marshal(flowDataResponse)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonStr = append(jsonStr, []byte("\n")...)
	w.Write(jsonStr)
}
