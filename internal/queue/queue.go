package queue

import "sync"

type FlowData struct {
	SrcApp  string `json:"src_app"`
	DestApp string `json:"dest_app"`
	VpcId   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytesRx int    `json:"bytes_rx"`
	Hour    int    `json:"hour"`
}

var fifoQueue []FlowData
var lock sync.Mutex

func Push(flowData []FlowData) {
	lock.Lock()
	defer lock.Unlock()

	fifoQueue = append(fifoQueue, flowData...)
}

func Consume() []FlowData {
	lock.Lock()
	defer lock.Unlock()

	flowData := fifoQueue
	fifoQueue = nil
	return flowData
}
