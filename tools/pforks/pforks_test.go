package pforks

import (
	"fmt"
	"testing"
)

func named() (code, message int) {
	return 1, 2
}

func TestFork_Run(t *testing.T) {
	f := &Fork{
		DoChildren: func() error {
			fmt.Println("this is child")
			return nil
		},
		Files:      nil,
	}

	f.DoMaster = func() error {
		fmt.Println("this is master")

		return nil
	}

	err := f.Run()
	if err != nil{
		panic(err)
	}
}
func TestName(t *testing.T) {

	fmt.Println(named())
}
