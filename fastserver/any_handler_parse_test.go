package fastserver

import (
	"fmt"
	"testing"
)

type C struct {
	A int
}
type Model struct {
	C
}

func Test_parseFunc(t *testing.T) {

	res, err := parseFuncHandler(func(ctx *Context, model *Model) (code int, resp *Model) {

		return 0, nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.response)
	fmt.Println(res.factory())
}
