package store

import "sync"

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

type flowDataStore struct {
	lock sync.RWMutex
	init bool
	kv   map[int]map[FlowDataKey]FlowDataValue
}

var store flowDataStore

func initialize() {
	store.kv = make(map[int]map[FlowDataKey]FlowDataValue)
	store.init = true
}

func Save(key FlowDataKey, value FlowDataValue) {
	store.lock.Lock()
	if store.init == false {
		initialize()
	}

	valueMap := store.kv[key.Hour]
	if valueMap == nil {
		store.kv[key.Hour] = make(map[FlowDataKey]FlowDataValue)
	}
	store.kv[key.Hour][key] = value

	store.lock.Unlock()
}

func LookupValue(key FlowDataKey) FlowDataValue {
	store.lock.Lock()
	if store.init == false {
		initialize()
	}
	store.lock.Unlock()

	store.lock.RLock()
	defer store.lock.RUnlock()

	return store.kv[key.Hour][key]
}

func LookupHour(hour int) map[FlowDataKey]FlowDataValue {
	store.lock.Lock()
	if store.init == false {
		initialize()
	}
	store.lock.Unlock()

	store.lock.RLock()
	defer store.lock.RUnlock()

	return store.kv[hour]
}
