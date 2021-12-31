package tools

import (
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"sync"
	"time"
)

type AddrClient struct {
	clients map[string]*fasthttp.Client
	lock sync.Mutex
}

func genKey(addr string,timeout int)string{
	return addr +strconv.Itoa(timeout)
}

func (c *AddrClient)GetClient(addr string,timeout int)*fasthttp.Client{
	c.lock.Lock()
	key := genKey(addr,timeout)
	cli := c.clients[key]
	if cli == nil{
		cli = &fasthttp.Client{
			ReadTimeout: time.Duration(timeout)*time.Second,
			Dial:  func(string) (net.Conn, error){
				return net.Dial("tcp",addr)
			},

		}
		c.clients[key] = cli
	}
	c.lock.Unlock()
	return cli
}

var clients = &AddrClient{clients: map[string]*fasthttp.Client{}}


func FastHttpClient(addr string,timeout int )*fasthttp.Client{
	return clients.GetClient(addr,timeout)
}
