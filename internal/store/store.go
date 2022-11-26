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

type FlowDataStore struct {
	lock sync.RWMutex
	init bool
	kv   map[int]map[FlowDataKey]FlowDataValue
}

func NewFlowDataStore() *FlowDataStore {
	return &FlowDataStore{
		lock: sync.RWMutex{},
		kv:   map[int]map[FlowDataKey]FlowDataValue{},
		init: false,
	}
}

func (s *FlowDataStore) Save(key FlowDataKey, value FlowDataValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	valueMap := s.kv[key.Hour]
	if valueMap == nil {
		s.kv[key.Hour] = make(map[FlowDataKey]FlowDataValue)
	}
	s.kv[key.Hour][key] = value
}

func (s *FlowDataStore) LookupValue(key FlowDataKey) FlowDataValue {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.kv[key.Hour][key]
}

func (s *FlowDataStore) LookupHour(hour int) map[FlowDataKey]FlowDataValue {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.kv[hour]
}
