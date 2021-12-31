package core

import (
	"fmt"
	"testing"
)

func Test_httpHealthyCheck_check(t *testing.T) {
	check := &httpHealthyCheck{
		Host:         "10.1.87.70",
		Method:       "GET",
		Path:         "/v2/iat",
		Body:         "",
		SuccessCodes: []int{200},
	}

	fmt.Println(check.check("10.1.87.70:8000"))
}
