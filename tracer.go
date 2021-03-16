package trace

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Stage store trace information in order during the same stage
// stageName => []key=>value
// for example,
// {
//	"gather_information": [
//      {
//			"gather_age": "success"
//		},
//		{
//			"gather_interest": "fail"
//		}
//	 ]
// }
type Stage map[string][]map[string]interface{}

// StageSlice stage slice
type StageSlice []Stage

type Tracer struct {
	sync.Mutex

	stages     StageSlice     // all stages in call order
	stageIndex map[string]int // stage => index
}

func NewTracer() *Tracer {
	return &Tracer{
		stageIndex: make(map[string]int),
	}
}

// AddKeyValue add key=>value trace information append to stage
func (t *Tracer) AddKeyValue(stage string, key string, data interface{}) {
	t.AddInfo(stage, map[string]interface{}{key: data})
}

// AddInfo add map trace information append to stage
func (t *Tracer) AddInfo(stageName string, info map[string]interface{}) {
	t.Lock()
	defer t.Unlock()

	index := t.getStageIndex(stageName)
	t.stages[index][stageName] = append(t.stages[index][stageName], info)
}

func (t *Tracer) getStageIndex(stageName string) int {
	if order, ok := t.stageIndex[stageName]; ok {
		return order
	} else {
		order := len(t.stageIndex)
		t.stageIndex[stageName] = order
		t.stages = append(t.stages, make(Stage))
		return order
	}
}

// Marshal
func (t *Tracer) Marshal() ([]byte, error) {
	if bytes, err := json.Marshal(t.stages); err != nil {
		return nil, err
	} else {
		return bytes, err
	}
}

// Unmarshal
func Unmarshal(data []byte) (*Tracer, error) {
	var stages []Stage

	if err := json.Unmarshal(data, &stages); err != nil {
		return nil, err
	}

	stageIndex := make(map[string]int)
	order := 0
	for _, stageInfo := range stages {
		for stage := range stageInfo {
			stageIndex[stage] = order
			order++
		}
	}

	return &Tracer{
		stages:     stages,
		stageIndex: stageIndex,
	}, nil
}

// Merge merge data unmarshal by another whiteboard
func (t *Tracer) Merge(data []byte) error {
	other, err := Unmarshal(data)
	if err != nil {
		return fmt.Errorf("string[%s] not valid white board, err: %s", string(data), err)
	}

	return t.MergeTracer(other)
}

// MergeTracer merge another tracer information
func (t *Tracer) MergeTracer(other *Tracer) error {
	for _, stageInfo := range other.stages {
		for stageName, dataSlice := range stageInfo {
			for _, data := range dataSlice {
				for key, value := range data {
					t.AddKeyValue(stageName, key, value)
				}
			}
		}
	}

	return nil
}
