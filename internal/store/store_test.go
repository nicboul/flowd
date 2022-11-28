package store

import (
	"fmt"
	"testing"
)

func Test_Store(t *testing.T) {

	tests := []struct {
		name            string
		tuple           *FlowDataTuple
		value           *FlowDataValue
		expectedByHour  map[FlowDataTuple]FlowDataValue
		expectedByTuple FlowDataValue
	}{
		{
			name:  "empty store",
			tuple: &FlowDataTuple{},
			value: &FlowDataValue{},
			expectedByHour: map[FlowDataTuple]FlowDataValue{
				// What is the behavior we want when we send "invalid" data to the store
				// what do we consider "invalid" ?
				{Hour: 0}: {0, 0},
			},
			expectedByTuple: FlowDataValue{},
		},
		{
			name: "with values",
			tuple: &FlowDataTuple{
				SrcApp:  "srcapp",
				DestApp: "destapp",
				VpcId:   "vpcid",
				Hour:    1,
			},
			value: &FlowDataValue{100, 200},
			expectedByHour: map[FlowDataTuple]FlowDataValue{
				{"srcapp", "destapp", "vpcid", 1}: {100, 200},
			},
			expectedByTuple: FlowDataValue{100, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewFlowDataStore()

			store.Save(tt.tuple, tt.value)

			byHour := store.LookupByHour(tt.tuple.Hour)
			byTuple := store.LookupByTuple(*tt.tuple)

			if !(fmt.Sprint(byTuple) == fmt.Sprint(tt.expectedByTuple)) {
				t.Fatalf("unexpected result:\n- want: %v\n- got: %v",
					tt.expectedByTuple, byTuple)
			}

			if !(fmt.Sprint(byHour) == fmt.Sprint(tt.expectedByHour)) {
				t.Fatalf("unexpected result:\n- want: %v\n- got: %v",
					tt.expectedByHour, byHour)
			}
		})
	}
}
