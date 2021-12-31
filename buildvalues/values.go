package buildvalues

import "log"

var(
	Mode = Debug
)

const (
	Debug = "debug"
	Version = "1.0.0"
)

func init(){
	log.Println("server run at mode:",Mode)
}

