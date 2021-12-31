package dns

import (
	"fmt"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/tools"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)
/*
Deprecated this package have not fully implemented the DNS protocol.
 */

func ServeDNS(port int) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("DNS Server Listing ... at ",port)
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)

		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go serverDNS(addr, conn, msg)
	}
}

// serverDNS serve
func serverDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	// query info
	if len(msg.Questions) < 1 {
		return
	}
	question := msg.Questions[0]
	var (
		queryNameStr = question.Name.String()
		queryType    = question.Type
		queryName, _ = dnsmessage.NewName(queryNameStr)
	)

	logger.Debug().Info("stdDns query:",queryType.String(),";",queryNameStr)

	var resource dnsmessage.Resource
	switch queryType {
	case dnsmessage.TypeA:
		rst, err := resolveDns(queryNameStr[:len(queryNameStr)-1], addr.IP.To4())
		if err != nil {
			response(addr, conn, msg)
			return
		}
		resource = NewAResource(queryName, rst)

	default:
		logger.Err().Error("not support dns queryType:",queryType.String())
		response(addr, conn, msg)
		return
	}

	// send response
	msg.Response = true
	msg.Answers = append(msg.Answers, resource)
	response(addr, conn, msg)
}

func resolveDns(host string, clientIp net.IP) ([4]byte, error) {
	ctx := &core.Context{}
	core.InitRequestContext(ctx, clientIp, "")
	var addr core.Address
	core.ResolveIpsByHost(ctx, host, &addr)
	if len(addr.Ips) > 0 {
		ip, _ := tools.SplitIp(addr.Ips[0])
		r := net.ParseIP(ip).To4()
		return [4]byte{r[0], r[1], r[2], r[3]}, nil
	}
	return [4]byte{}, fmt.Errorf("not found")

}

// Response return
func response(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	packed, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(packed, addr); err != nil {
		//fmt.Println(err)
		logger.Err().Error("write dns error:%w",err)
	}
}

// NewAResource A record
func NewAResource(query dnsmessage.Name, a [4]byte) dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
			TTL:   600,
		},
		Body: &dnsmessage.AResource{
			A: a,
		},
	}
}

// NewPTRResource PTR record
func NewPTRResource(query dnsmessage.Name, ptr string) dnsmessage.Resource {
	name, _ := dnsmessage.NewName(ptr)
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
		},
		Body: &dnsmessage.PTRResource{
			PTR: name,
		},
	}
}
