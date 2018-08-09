package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMapKeyWithSpace(t *testing.T) {

	type response2 struct {
		Page   int      `json:"page"`
		Fruits []string `json:"fruits"`
	}

	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))



}
