package api

import (
	"fmt"
	"testing"
)

type a struct {
	Ass fmt.Stringer
}

type imp struct {

}

func (i imp) String() string {
	return "haha"
}

func Test_injector_doInject(t *testing.T) {
	var s = &a{}
	fmt.Println(s.Ass.String())
}
