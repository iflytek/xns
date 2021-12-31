package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"strconv"
)

func main(){
	sti := false
	flag.BoolVar(&sti,"sti",sti,"string to int")
	flag.Parse()
	if sti{
		for _, s := range flag.Args(){
			fmt.Println(ip2int(s))
		}
	}else{
		for _, s := range flag.Args() {
			i,_:=strconv.Atoi(s)
			fmt.Println(int2ip(uint32(i)))
		}
	}

}

func ip2int(ip string)uint32{
	return  binary.BigEndian.Uint32(net.ParseIP(ip).To4())
}

func int2ip(i uint32)string{
	b:= make([]byte,4)
	binary.BigEndian.PutUint32(b,i)
	return fmt.Sprintf("%d.%d.%d.%d",b[0],b[1],b[2],b[3])
}
