package serial

import (
	"bytes"
	"encoding/gob"
	"fmt"
	reflect "reflect"
	"sync"
	"testing"

	"ergo.services/ergo/gen"
	"ergo.services/ergo/lib"
	"ergo.services/ergo/net/edf"
	"google.golang.org/protobuf/proto"
)

var (
	cache           = &sync.Map{}
	regCache        = &sync.Map{}
	edfOptionsCache = edf.Options{
		Cache:    cache,
		RegCache: regCache,
	}
)

func init() {
	edf.RegisterTypeOf(struct {
		A float32
		B float64
		C string
	}{})
	edf.RegisterTypeOf(ComplexStructValue{})
	edf.RegisterTypeOf(NestedStructValue{})
	edf.RegisterTypeOf(SimpleStructValue{})

	regTypeName := func(t reflect.Type) string {
		return fmt.Sprintf("#%s/%s", t.PkgPath(), t.Name())
	}

	fmt.Println(regTypeName(reflect.TypeOf(ComplexStructValue{})))
	regCache.Store(reflect.TypeOf(ComplexStructValue{}), []byte{131, 0x13, 0x88})
	regCache.Store(uint16(5000), "#serial/ComplexStructValue")
	regCache.Store(reflect.TypeOf(NestedStructValue{}), []byte{131, 0x13, 0x89})
	regCache.Store(uint16(5001), "#serial/NestedStructValue")
	regCache.Store(reflect.TypeOf(SimpleStructValue{}), []byte{131, 0x13, 0x8a})
	regCache.Store(uint16(5002), "#serial/SimpleStructValue")
}

// =============================================================================
// TYPE CACHING BENCHMARKS - EDF Optimized with Type Cache
// =============================================================================
// These benchmarks demonstrate the performance benefits of using EDF's
// type caching feature which significantly reduces memory allocations
// and improves encoding/decoding speed for complex types.

// Benchmark Encode String with Type Cache
func BenchmarkEncodeStringCached(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode("Ergo Framework", buf, edfOptionsCache); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode String with Type Cache
func BenchmarkDecodeStringCached(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	options := edf.Options{Cache: cache}

	if err := edf.Encode("Ergo Framework", buf, edfOptionsCache); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode Map with Type Cache (should show biggest improvement)
func BenchmarkEncodeMapCached(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := map[string]SimpleStructValue{
		"key1": {Name: "value1", Id: 1},
		"key2": {Name: "value2", Id: 2},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode Map with Type Cache
func BenchmarkDecodeMapCached(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := map[string]SimpleStructValue{
		"key1": {Name: "value1", Id: 1},
		"key2": {Name: "value2", Id: 2},
	}

	if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edfOptionsCache)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode Complex Struct with Type Cache
func BenchmarkEncodeComplexStructCached(b *testing.B) {
	value := ComplexStructValue{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
		ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode Complex Struct with Type Cache
func BenchmarkDecodeComplexStructCached(b *testing.B) {
	value := ComplexStructValue{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
		ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edfOptionsCache)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode Nested Struct with Type Cache (should show massive improvement)
func BenchmarkEncodeNestedStructCached(b *testing.B) {
	value := NestedStructValue{
		Name: "nested_test",
		Id:   456,
		Complex: ComplexStructValue{
			Name:      "complex_value",
			Id:        789,
			Tags:      []string{"nested_tag1", "nested_tag2"},
			Metadata:  map[string]string{"nested_key": "nested_value"},
			Pid:       gen.PID{Node: "nested_node", ID: 2, Creation: 2},
			ProcessId: gen.ProcessID{Node: "nested_node", Name: "nested_process"},
		},
		ComplexMap: map[string]ComplexStructValue{
			"complex1": {
				Name:     "complex1",
				Id:       100,
				Tags:     []string{"tag1"},
				Metadata: map[string]string{"key1": "value1"},
			},
			"complex2": {
				Name:     "complex2",
				Id:       200,
				Tags:     []string{"tag2"},
				Metadata: map[string]string{"key2": "value2"},
			},
		},
		NestedMap: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode Nested Struct with Type Cache
func BenchmarkDecodeNestedStructCached(b *testing.B) {
	value := NestedStructValue{
		Name: "nested_test",
		Id:   456,
		Complex: ComplexStructValue{
			Name:      "complex_value",
			Id:        789,
			Tags:      []string{"nested_tag1", "nested_tag2"},
			Metadata:  map[string]string{"nested_key": "nested_value"},
			Pid:       gen.PID{Node: "nested_node", ID: 2, Creation: 2},
			ProcessId: gen.ProcessID{Node: "nested_node", Name: "nested_process"},
		},
		ComplexMap: map[string]ComplexStructValue{
			"complex1": {
				Name:     "complex1",
				Id:       100,
				Tags:     []string{"tag1"},
				Metadata: map[string]string{"key1": "value1"},
			},
			"complex2": {
				Name:     "complex2",
				Id:       200,
				Tags:     []string{"tag2"},
				Metadata: map[string]string{"key2": "value2"},
			},
		},
		NestedMap: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	if err := edf.Encode(value, buf, edfOptionsCache); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edfOptionsCache)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// ORIGINAL BENCHMARKS (for comparison)
// =============================================================================

// Benchmark Encode String
func BenchmarkEncodeString(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := edf.Encode("Ergo Framework", buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeStringProtobuf(b *testing.B) {
	value := &TestStruct{
		A: 123.45,
		B: 678.90,
		C: []byte("Ergo Framework"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeStringGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := "Ergo Framework"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode PID
func BenchmarkEncodePID(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := gen.PID{Node: "demo@127.0.0.1", ID: 312, Creation: 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodePIDProtobuf(b *testing.B) {
	value := &PID{
		Node:     "demo@127.0.0.1",
		Id:       312,
		Creation: 2,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodePIDGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := gen.PID{Node: "demo@127.0.0.1", ID: 312, Creation: 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode ProcessID
func BenchmarkEncodeProcessID(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := gen.ProcessID{Node: "demo@127.0.0.1", Name: "example"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeProcessIDProtobuf(b *testing.B) {
	value := &ProcessID{
		Node: "demo@127.0.0.1",
		Name: "example",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeProcessIDGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := gen.ProcessID{Node: "demo@127.0.0.1", Name: "example"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Encode Struct
func BenchmarkEncodeStruct(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := struct {
		A float32
		B float64
		C string
	}{
		A: 123.45,
		B: 678.90,
		C: "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeStructProtobuf(b *testing.B) {
	value := &TestStruct{
		A: 123.45,
		B: 678.90,
		C: []byte("test"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeStructGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := struct {
		A float32
		B float64
		C string
	}{
		A: 123.45,
		B: 678.90,
		C: "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode String
func BenchmarkDecodeString(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	if err := edf.Encode("Ergo Framework", buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeStringProtobuf(b *testing.B) {
	value := &TestStruct{
		A: 123.45,
		B: 678.90,
		C: []byte("Ergo Framework"),
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := proto.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeStringGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := "Ergo Framework"
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result string
		dec := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode PID
func BenchmarkDecodePID(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := gen.PID{Node: "demo@127.0.0.1", ID: 312, Creation: 2}
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodePIDProtobuf(b *testing.B) {
	value := &PID{
		Node:     "demo@127.0.0.1",
		Id:       312,
		Creation: 2,
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result PID
		err := proto.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodePIDGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := gen.PID{Node: "demo@127.0.0.1", ID: 312, Creation: 2}
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result gen.PID
		dec := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode ProcessID
func BenchmarkDecodeProcessID(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := gen.ProcessID{Node: "demo@127.0.0.1", Name: "example"}
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeProcessIDProtobuf(b *testing.B) {
	value := &ProcessID{
		Node: "demo@127.0.0.1",
		Name: "example",
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result ProcessID
		err := proto.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeProcessIDGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := gen.ProcessID{Node: "demo@127.0.0.1", Name: "example"}
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result gen.ProcessID
		dec := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Decode Struct
func BenchmarkDecodeStruct(b *testing.B) {
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)

	value := struct {
		A float32
		B float64
		C string
	}{
		A: 123.45,
		B: 678.90,
		C: "test",
	}
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := edf.Decode(buf.B, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeStructProtobuf(b *testing.B) {
	value := &TestStruct{
		A: 123.45,
		B: 678.90,
		C: []byte("test"),
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := proto.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeStructGob(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	value := struct {
		A float32
		B float64
		C string
	}{
		A: 123.45,
		B: 678.90,
		C: "test",
	}
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result struct {
			A float32
			B float64
			C string
		}
		dec := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

type ComplexStructValue struct {
	Name      string
	Id        int32
	Tags      []string
	Metadata  map[string]string
	Pid       gen.PID
	ProcessId gen.ProcessID
}

type NestedStructValue struct {
	Name       string
	Id         int32
	Complex    ComplexStructValue
	ComplexMap map[string]ComplexStructValue
	NestedMap  map[string]string
}

func BenchmarkEncodeComplexStructEDF(b *testing.B) {
	value := ComplexStructValue{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
		ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeComplexStructEDF(b *testing.B) {
	value := ComplexStructValue{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
		ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}
	data := buf.B
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _, err := edf.Decode(data, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
		_ = res.(ComplexStructValue)
	}
}

func BenchmarkEncodeComplexStructProtobuf(b *testing.B) {
	value := &ComplexStruct{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
		ProcessId: &ProcessID{Node: "node1", Name: "process1"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeComplexStructProtobuf(b *testing.B) {
	value := &ComplexStruct{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
		ProcessId: &ProcessID{Node: "node1", Name: "process1"},
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result ComplexStruct
		if err := proto.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeComplexStructGob(b *testing.B) {
	value := &ComplexStruct{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
		ProcessId: &ProcessID{Node: "node1", Name: "process1"},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeComplexStructGob(b *testing.B) {
	value := &ComplexStruct{
		Name:      "test",
		Id:        123,
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
		ProcessId: &ProcessID{Node: "node1", Name: "process1"},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result ComplexStruct
		dec := gob.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeNestedStructEDF(b *testing.B) {
	value := NestedStructValue{
		Name: "test",
		Id:   123,
		Complex: ComplexStructValue{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
			ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]ComplexStructValue{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
				ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedStructEDF(b *testing.B) {
	value := NestedStructValue{
		Name: "test",
		Id:   123,
		Complex: ComplexStructValue{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
			ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]ComplexStructValue{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       gen.PID{Node: "node1", ID: 1, Creation: 1},
				ProcessId: gen.ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}
	data := buf.B
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _, err := edf.Decode(data, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
		_ = res.(NestedStructValue)
	}
}

func BenchmarkEncodeNestedStructProtobuf(b *testing.B) {
	value := &NestedStruct{
		Name: "test",
		Id:   123,
		Complex: &ComplexStruct{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
			ProcessId: &ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]*ComplexStruct{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
				ProcessId: &ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedStructProtobuf(b *testing.B) {
	value := &NestedStruct{
		Name: "test",
		Id:   123,
		Complex: &ComplexStruct{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
			ProcessId: &ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]*ComplexStruct{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
				ProcessId: &ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result NestedStruct
		if err := proto.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeNestedStructGob(b *testing.B) {
	value := &NestedStruct{
		Name: "test",
		Id:   123,
		Complex: &ComplexStruct{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
			ProcessId: &ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]*ComplexStruct{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
				ProcessId: &ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedStructGob(b *testing.B) {
	value := &NestedStruct{
		Name: "test",
		Id:   123,
		Complex: &ComplexStruct{
			Name:      "test",
			Id:        123,
			Tags:      []string{"tag1", "tag2"},
			Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
			Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
			ProcessId: &ProcessID{Node: "node1", Name: "process1"},
		},
		ComplexMap: map[string]*ComplexStruct{
			"key1": {
				Name:      "test",
				Id:        123,
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
				Pid:       &PID{Node: "node1", Id: 1, Creation: 1},
				ProcessId: &ProcessID{Node: "node1", Name: "process1"},
			},
		},
		NestedMap: map[string]string{"key1": "value1", "key2": "value2"},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result NestedStruct
		dec := gob.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

type SimpleStructValue struct {
	Name string
	Id   int32
}

func BenchmarkEncodeMapEDF(b *testing.B) {
	value := map[string]SimpleStructValue{
		"key1": {Name: "test1", Id: 123},
		"key2": {Name: "test2", Id: 456},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMapEDF(b *testing.B) {
	value := map[string]SimpleStructValue{
		"key1": {Name: "test1", Id: 123},
		"key2": {Name: "test2", Id: 456},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}
	data := buf.B
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _, err := edf.Decode(data, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
		_ = res.(map[string]SimpleStructValue)
	}
}

func BenchmarkEncodeMapProtobuf(b *testing.B) {
	value := &MapMessage{
		Map: map[string]*SimpleStruct{
			"key1": {Name: "test1", Id: 123},
			"key2": {Name: "test2", Id: 456},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMapProtobuf(b *testing.B) {
	value := &MapMessage{
		Map: map[string]*SimpleStruct{
			"key1": {Name: "test1", Id: 123},
			"key2": {Name: "test2", Id: 456},
		},
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result MapMessage
		if err := proto.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeMapGob(b *testing.B) {
	value := map[string]*SimpleStruct{
		"key1": {Name: "test1", Id: 123},
		"key2": {Name: "test2", Id: 456},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMapGob(b *testing.B) {
	value := map[string]*SimpleStruct{
		"key1": {Name: "test1", Id: 123},
		"key2": {Name: "test2", Id: 456},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]*SimpleStruct
		dec := gob.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeNestedMapEDF(b *testing.B) {
	value := map[string]map[string]SimpleStructValue{
		"outer1": {
			"inner1": {Name: "test1", Id: 123},
			"inner2": {Name: "test2", Id: 456},
		},
		"outer2": {
			"inner3": {Name: "test3", Id: 789},
		},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := edf.Encode(value, buf, edf.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedMapEDF(b *testing.B) {
	value := map[string]map[string]SimpleStructValue{
		"outer1": {
			"inner1": {Name: "test1", Id: 123},
			"inner2": {Name: "test2", Id: 456},
		},
		"outer2": {
			"inner3": {Name: "test3", Id: 789},
		},
	}
	buf := lib.TakeBuffer()
	defer lib.ReleaseBuffer(buf)
	if err := edf.Encode(value, buf, edf.Options{}); err != nil {
		b.Fatal(err)
	}
	data := buf.B
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _, err := edf.Decode(data, edf.Options{})
		if err != nil {
			b.Fatal(err)
		}
		_ = res.(map[string]map[string]SimpleStructValue)
	}
}

func BenchmarkEncodeNestedMapProtobuf(b *testing.B) {
	value := &NestedMapMessage{
		Map: map[string]*MapMessage{
			"outer1": {Map: map[string]*SimpleStruct{
				"inner1": {Name: "test1", Id: 123},
				"inner2": {Name: "test2", Id: 456},
			}},
			"outer2": {Map: map[string]*SimpleStruct{
				"inner3": {Name: "test3", Id: 789},
			}},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedMapProtobuf(b *testing.B) {
	value := &NestedMapMessage{
		Map: map[string]*MapMessage{
			"outer1": {Map: map[string]*SimpleStruct{
				"inner1": {Name: "test1", Id: 123},
				"inner2": {Name: "test2", Id: 456},
			}},
			"outer2": {Map: map[string]*SimpleStruct{
				"inner3": {Name: "test3", Id: 789},
			}},
		},
	}
	data, err := proto.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result NestedMapMessage
		if err := proto.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeNestedMapGob(b *testing.B) {
	value := map[string]map[string]*SimpleStruct{
		"outer1": {
			"inner1": {Name: "test1", Id: 123},
			"inner2": {Name: "test2", Id: 456},
		},
		"outer2": {
			"inner3": {Name: "test3", Id: 789},
		},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedMapGob(b *testing.B) {
	value := map[string]map[string]*SimpleStruct{
		"outer1": {
			"inner1": {Name: "test1", Id: 123},
			"inner2": {Name: "test2", Id: 456},
		},
		"outer2": {
			"inner3": {Name: "test3", Id: 789},
		},
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]map[string]*SimpleStruct
		dec := gob.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&result); err != nil {
			b.Fatal(err)
		}
	}
}
