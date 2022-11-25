package store

type FlowDataKey struct {
	SrcApp  string
	DestApp string
	VpcId   string
	Hour    int
}

type FlowDataValue struct {
	BytesTx int
	BytesRx int
}

var initStore bool = false
var store map[int]map[FlowDataKey]FlowDataValue

func initialize() {
	store = make(map[int]map[FlowDataKey]FlowDataValue)
	initStore = true
}

func Save(key FlowDataKey, value FlowDataValue) {
	if initStore == false {
		initialize()
	}

	valueMap := store[key.Hour]
	if valueMap == nil {
		store[key.Hour] = make(map[FlowDataKey]FlowDataValue)
	}
	store[key.Hour][key] = value
}

func LookupValue(key FlowDataKey) FlowDataValue {
	if initStore == false {
		initialize()
	}

	return store[key.Hour][key]
}

func LookupHour(hour int) map[FlowDataKey]FlowDataValue {
	return store[hour]
}
