## Trace
The trace library can help you inspect what happens during a call. All the trace information is stored and marshaled
in call order.

With multiple-stage support, you can gather related trace information in the same stage, which helps group and
classify information to achieve better readability after marshaled.


## Usage
```
tracer := NewTracer()

tracer.AddKeyValue("stage1", "key", interface{})
tracer.AddKeyValue("stage2", "key", interface{})

fmt.Println(tracer.Marshal())
```

