package queue

type FlowData struct {
	SrcApp  string `json:"src_app"`
	DestApp string `json:"dest_app"`
	VpcId   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytesRx int    `json:"bytes_rx"`
	Hour    int    `json:"hour"`
}

type FlowDataQueue struct {
	Channel chan []FlowData
}

func NewFlowDataQueue(size int) *FlowDataQueue {
	return &FlowDataQueue{
		Channel: make(chan []FlowData, size),
	}
}

func (q *FlowDataQueue) TryEnqueue(flowData []FlowData) bool {
	select {
	case q.Channel <- flowData:
		return true
	default:
		return false
	}
}
