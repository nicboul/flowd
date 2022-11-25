package flowdata

type FlowData struct {
	SrcApp  string `json:"src_app"`
	DestApp string `json:"dest_app"`
	VpcId   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytesRx int    `json:"bytes_rx"`
	Hour    int    `json:"hour"`
}
