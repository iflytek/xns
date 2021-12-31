package msc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)


func TestDnsResponse_Json(t *testing.T) {


}

func BenchmarkJson(b *testing.B) {
	bf := bytes.Buffer{}
	bf.Grow(200)


}

func TestJson2(t *testing.T) {
	s:= Sip{
		Compel: &Compel{Value: false},
	}
	v,_:=json.Marshal(s)

	fmt.Println(string(v))
}
