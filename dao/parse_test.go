package dao

import (
	"fmt"
	"reflect"
	"testing"
)
type poolA struct {
	Description    string                 `json:"description"`
	Name           string                 `json:"name" format:"name"`
	LbMode         int                    `json:"lb_mode" desc:"负载均衡模式"`                         // 1，就近分配原则，并且根据权重来分配 2，某一地机房全部没了，更具权重分配到其他机房。
	LbConfig       map[string]interface{} `json:"lb_config" desc:"负载均衡配置"`                       // 负载均衡配置
	FailOverConfig map[string]interface{} `json:"fail_over_config" desc:"负载均衡选择失败时，兜底配置,object"` // fail配置 ，dx:[hu,gz],hu:[dx,gz],gz:[]
}

func Test_parseUpdateTag(t *testing.T) {
	fmt.Println(parseUpdateTag(reflect.ValueOf(poolA{}), func(s string) string {
		return s
	}))

}
