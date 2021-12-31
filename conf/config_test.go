package conf

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/bytebufferpool"
	"testing"
)

func Test_parseCfgFromFile(t *testing.T) {
	fmt.Println(parseCfgFromFile(`./nameserver.cfg`))
	bs,_:=json.Marshal(Load())
	fmt.Println(string(bs))
}

func TestBD(t *testing.T) {
	bytebufferpool.Get()

}
