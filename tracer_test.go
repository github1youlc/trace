package trace

import (
	"testing"
)

type DebugData struct {
	A int
	B string
}

func TestTracer(t *testing.T) {
	tracer := NewTracer()

	tracer.AddKeyValue("stage1", "req", "stage1_req")
	tracer.AddKeyValue("stage2", "req", "stage2_req")
	tracer.AddKeyValue("stage1", "resp", "stage1_resp")
	tracer.AddKeyValue("stage2", "resp", "stage2_resp")

	data, err := tracer.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

	tracer2, err := Unmarshal(data)
	if err != nil {
		t.Error(err)
	}

	tracer2.AddKeyValue("stage3", "step1", map[string]interface{}{"data1": 1, "date2": "test"})
	tracer2.AddKeyValue("stage3", "step2", &DebugData{
		A: 1,
		B: "xxx",
	})
	data2, err := tracer2.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data2))

	err = tracer.Merge(data2)
	if err != nil {
		t.Error(err)
	}

	data3, err := tracer.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data3))
}
