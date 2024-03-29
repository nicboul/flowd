package flowdatawrite

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nicboul/flowdata/internal/queue"
)

type FlowDataWrite struct {
	Queue *queue.FlowDataQueue
}

func NewFlowDataWrite(queue *queue.FlowDataQueue) *FlowDataWrite {
	return &FlowDataWrite{
		Queue: queue,
	}
}

func (f *FlowDataWrite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var flowData []queue.FlowData
	err = json.Unmarshal(b, &flowData)
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If we fail to queue the message, we need to notice the caller
	success := f.Queue.TryEnqueue(flowData)
	if !success {
		http.Error(w, "reached max capacity", http.StatusInternalServerError)
		return
	}
}
