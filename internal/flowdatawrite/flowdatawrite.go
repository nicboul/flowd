package flowdatawrite

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nicboul/flowdata/internal/aggregator"
	"github.com/nicboul/flowdata/internal/flowdata"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
)

type FlowDataWrite struct {
	Store *store.FlowDataStore
}

func (f *FlowDataWrite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var flowData []flowdata.FlowData
	err = json.Unmarshal(b, &flowData)
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queue.Push(flowData)

	aggregator.Aggregate(f.Store)
}
