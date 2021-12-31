package main

import (
	"net"
	"testing"
)

func TestV6(t *testing.T) {

	conn ,err := net.Dial("tcp","[2400:da00::dbf:0:100]:80")
	if err != nil{
		panic(err)
	}
	conn.Close()

}
