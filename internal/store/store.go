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
	Lock sync.RWMutex
	init bool
	kv   map[int]map[FlowDataTuple]FlowDataValue
}

func NewFlowDataStore() *FlowDataStore {
	return &FlowDataStore{
		Lock: sync.RWMutex{},
		kv:   map[int]map[FlowDataTuple]FlowDataValue{},
		init: false,
	}
}

func (s *FlowDataStore) SaveWithLock(key *FlowDataTuple, value *FlowDataValue) error {
	valueMap := s.kv[key.Hour]
	if valueMap == nil {
		s.kv[key.Hour] = make(map[FlowDataTuple]FlowDataValue)
	}
	s.kv[key.Hour][*key] = *value

	return nil
}

func (s *FlowDataStore) LookupByTupleWithLock(key FlowDataTuple) FlowDataValue {
	return s.kv[key.Hour][key]
}

func (s *FlowDataStore) LookupByHour(hour int) map[FlowDataTuple]FlowDataValue {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	return s.kv[hour]
}
