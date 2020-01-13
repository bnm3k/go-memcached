package lfu

import (
	"reflect"
	"testing"
)

func TestLRUCacheLeetCode(t *testing.T) {
	type testInput struct {
		key     string
		val     string
		exptime int
	}
	lfuCache := Constructor(2)
	actions := []string{"put", "put", "get", "put", "get", "get", "put", "get", "get", "get"}
	inputs := []testInput{{"1", "1", 0}, {"2", "2", 0}, {"1", "", 0}, {"3", "3", 0},
		{"2", "", 0}, {"3", "", 0}, {"4", "4", 0}, {"1", "", 0}, {"3", "", 0}, {"4", "", 0}}
	expected := []string{"null", "null", "1", "null", "", "3", "null", "", "3", "4"}
	output := make([]string, len(actions))
	for i, action := range actions {
		input := inputs[i]
		switch action {
		case "put":
			lfuCache.Set(input.key, input.val, input.exptime)
			output[i] = "null"
		case "get":
			val, _ := lfuCache.Get(input.key)
			output[i] = val
		}
	}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("\ngot %v \nwant %v\n", output, expected)
	}
}
