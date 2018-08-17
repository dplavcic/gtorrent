package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestParseInt(t *testing.T) {
	var v int64
	data := []byte("2e")

	// reflect.ValueOf(v)  - not addressable
	// reflect.ValueOf(&x).Elem().Set()

	pv := reflect.ValueOf(&v)
	e := Decoder{buf: bytes.NewBuffer(data), v: v}
	e.readInt(pv.Elem())
	if 2 != v {
		t.Errorf("expected: %v, got %v", 2, v)
	}
}

func TestParseString(t *testing.T) {
	var s string
	data := []byte("3:foo")

	pv := reflect.ValueOf(&s)
	e := Decoder{buf: bytes.NewBuffer(data), v: s}
	e.readString(pv.Elem())
	if "foo" != s {
		t.Errorf("expected: %v, got %v", "foo", s)
	}
}

func TestParseList(t *testing.T) {
	var l []interface{}

	fmt.Printf("L with: %v\n", l)
	fmt.Printf("L with valueOf(): %v\n", reflect.ValueOf(l))
	fmt.Printf("L with type: %v\n", reflect.TypeOf(l))
	fmt.Printf("L with kind: %v\n", reflect.ValueOf(l).Kind())
	// l = make([]interface{}, 1, 1)
	// l[0] = "hello"
	data := []byte("3:fool3:barl3:dvae3:bazl1:ale1:bee4:tttte")
	// data := []byte("3:fooe")

	pv := reflect.ValueOf(&l)
	e := Decoder{buf: bytes.NewBuffer(data), v: l}
	e.readList(pv)
	t.Errorf("Result: %v\n", l)
}

type Person struct {
	Id    int8
	Name  P2
	First []byte
	Last  []interface{}
}
type P2 struct {
	Id    int
	Phone P3
}

type P3 struct {
	Country string
	Phone   string
}

func TestParseDict(t *testing.T) {
	p := Person{}

	data := []byte("2:Idi10e5:First3:foo4:Lastli1ei2ee4:Named2:Idi5e5:Phoned7:Country2:hr5:Phone3:385eee")
	pv := reflect.ValueOf(&p)
	e := Decoder{buf: bytes.NewBuffer(data), v: p}
	e.readDict(pv)
	t.Errorf("Result: %v\n", p)
}

func BenchmarkList(b *testing.B) {
	data := []byte("2:Idi10e5:First3:foo4:Lastli1ei2ee4:Named2:Idi5e5:Phoned7:Country2:hr5:Phone3:385eee")

	for n := 0; n < 10000; n++ {
		p := Person{}

		pv := reflect.ValueOf(&p)
		e := Decoder{buf: bytes.NewBuffer(data), v: p}
		e.readDict(pv)
	}
}

// BenchmarkList-8   	       1	1594964541 ns/op	41136712 B/op	 2210069 allocs/op
// PASS
// ok  	github.com/dplavcic/gtorrent/bencode	1.597s

// BenchmarkList-8         2000000000               0.02 ns/op            0 B/op          0 allocs/op
// PASS
// ok      github.com/dplavcic/gtorrent/bencode    0.249s
