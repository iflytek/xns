package dao

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_newCond(t *testing.T) {
	fmt.Println(newCond().eq("id","5555").and().eq("name","ffff").or().gt("age",5).contains("idc","ghj").String())
}

func BenchmarkJSON(b *testing.B) {
	var a  = map[string]interface{}{
		"Name":"xxx",
		"res":"xxx",
	}
	for i := 0; i < b.N; i++ {
		json.Marshal(a)
	}
}

