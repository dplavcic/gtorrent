package bencode

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestIntParser(t *testing.T) {
	testCases := map[string]struct {
		in  *bytes.Buffer
		out int64
	}{
		"0":          {bytes.NewBuffer([]byte("i0e")), 0},
		"100":        {bytes.NewBuffer([]byte("i100e")), 100},
		"-100":       {bytes.NewBuffer([]byte("i-100e")), -100},
		"2147483647": {bytes.NewBuffer([]byte("i2147483647e")), math.MaxInt32},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ReadInt(tt.in)
			if result != tt.out {
				t.Errorf("want: %d, got: %d\n", tt.out, result)
			}
		})
	}
}

func TestStringParser(t *testing.T) {
	testCases := map[string]struct {
		in  *bytes.Buffer
		out string
	}{
		"1:f":      {bytes.NewBuffer([]byte("1:f")), "f"},
		"2:fo":     {bytes.NewBuffer([]byte("2:fo")), "fo"},
		"6:f1o2o3": {bytes.NewBuffer([]byte("6:f1o2o3")), "f1o2o3"},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ReadString(tt.in)
			if result != tt.out {
				t.Errorf("want: %s, got: %s\n", tt.out, result)
			}
		})
	}
}

func TestListParser(t *testing.T) {
	testCases := map[string]struct {
		in  *bytes.Buffer
		out []interface{}
	}{
		"[[elem1]]":                {bytes.NewBuffer([]byte("ll3:fooee")), []interface{}{[]interface{}{"foo"}}},
		"[[elem1][elem2]]":         {bytes.NewBuffer([]byte("ll3:fooel3:baree")), []interface{}{[]interface{}{"foo"}, []interface{}{"bar"}}},
		"[[elem1, elem2] [elem3]]": {bytes.NewBuffer([]byte("ll3:foo3:barel3:bazee")), []interface{}{[]interface{}{"foo", "bar"}, []interface{}{"baz"}}},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ReadList(tt.in)
			tempResult := fmt.Sprintf("%s", result)
			tempOut := fmt.Sprint(tt.out)

			if bytes.Compare([]byte(tempResult), []byte(tempOut)) != 0 {
				t.Errorf("want: \n%s, got: \n%s\n", tt.out, result)
			}
		})
	}
}
