package store

import "sync"

type FlowDataTuple struct {
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
	kv   map[int]map[FlowDataTuple]FlowDataValue
}

type StoreManager interface {
	Save(FlowDataTuple, FlowDataValue) error
	LookupByFlowData(FlowDataTuple) FlowDataValue
	LookupByHour(int) map[FlowDataTuple]FlowDataValue
}

func NewFlowDataStore() *FlowDataStore {
	return &FlowDataStore{
		lock: sync.RWMutex{},
		kv:   map[int]map[FlowDataTuple]FlowDataValue{},
		init: false,
	}
}

func (s *FlowDataStore) Save(key *FlowDataTuple, value *FlowDataValue) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	valueMap := s.kv[key.Hour]
	if valueMap == nil {
		s.kv[key.Hour] = make(map[FlowDataTuple]FlowDataValue)
	}
	s.kv[key.Hour][*key] = *value

	return nil
}

func (s *FlowDataStore) LookupByTuple(key FlowDataTuple) FlowDataValue {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.kv[key.Hour][key]
}

func (s *FlowDataStore) LookupByHour(hour int) map[FlowDataTuple]FlowDataValue {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.kv[hour]
}
