package core

import (
	"encoding/json"
	"fmt"
	"github.com/xfyun/xns/tools"
	"github.com/xfyun/xns/tools/str"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

type healthyCheckFunc interface {
	check(addr string) error
}

type healthyCheckFactory func(name string, cfg string) (healthyCheckFunc, error)

var healthyCheckFactories = map[string]healthyCheckFactory{
	"tcp": func(name string, cfg string) (healthyCheckFunc, error) {
		check := &tcpHealthyCheck{}
		if cfg == "" {
			return check, nil
		}
		err := json.Unmarshal([]byte(cfg), check)
		if err != nil {
			return nil, fmt.Errorf("parse tcp healthy check:'%s' error:%s ", name, err.Error())
		}
		return check, nil
	},
	"http": func(name string, cfg string) (healthyCheckFunc, error) {
		check := &httpHealthyCheck{}
		if cfg == "" {
			return nil, fmt.Errorf("create http healthy error, config should not be empty")
		}
		err := json.Unmarshal([]byte(cfg), check)
		if err != nil {
			return nil, fmt.Errorf("parse tcp healthy check:'%s' error:%s ", name, err.Error())
		}
		if len(check.SuccessCodes) == 0 {
			check.SuccessCodes = []int{200, 201, 204, 302, 404}
		}
		if check.Timeout <= 0 {
			check.Timeout = 4
		}
		return check, nil
	},
}

func NewHealthyCheck(name string, cfg string) (healthyCheckFunc, error) {
	hf := healthyCheckFactories[name]
	if hf == nil {
		return nil, fmt.Errorf("create healthy func error,unknow factory name :%s", name)
	}
	return hf(name, cfg)
}

//tcp healthy check
type tcpHealthyCheck struct {
	Timeout int `json:"timeout"`
}

func (t *tcpHealthyCheck) check(addr string) error {

	timeout := t.Timeout
	if timeout <= 0 {
		timeout = 4
	}
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

type httpHealthyCheck struct {
	Host         string `json:"host"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Body         string `json:"body"`
	SuccessCodes []int  `json:"success_codes"`
	Timeout      int    `json:"timeout"`
}

func (h httpHealthyCheck) check(addr string) error {

	host := h.Host
	if host == ""{
		host = addr
	}
	request := fasthttp.AcquireRequest()
	request.SetHost(addr)
	request.Header.SetMethod(h.Method)
	request.URI().SetPath(h.Path)
	request.SetBody(str.BytesOf(h.Body))
	defer fasthttp.ReleaseRequest(request)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := tools.FastHttpClient(addr, h.Timeout).Do(request, resp)
	if err != nil {
		return err
	}

	bd := resp.Body()
	for _, code := range h.SuccessCodes {
		if code == resp.StatusCode() {
			return nil
		}
	}
	return fmt.Errorf("error code not expected :%d ,body:%s", resp.StatusCode(), string(bd))
}
