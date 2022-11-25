package queue

import (
	"github.com/nicboul/flowdata/internal/flowdata"
)

var fifoQueue []flowdata.FlowData

func Push(flowData []flowdata.FlowData) {
	fifoQueue = append(fifoQueue, flowData...)
}

func Consume() []flowdata.FlowData {
	flowData := fifoQueue
	fifoQueue = nil
	return flowData
}
