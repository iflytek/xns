package fastserver

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/prefork"
	"github.com/valyala/fasthttp/reuseport"
	"github.com/xfyun/xns/tools/str"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var stringOf = str.StringOf

func (c HandlersChain) Last() Handler {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// http server
// 基于fastHttp 封装路由功能
// 路由部分参考gin 的路由树，支持path 参数
type Server struct {
	RouterGroup     //
	closed          bool
	stopWg          sync.WaitGroup // 用于优雅启动停止时，等待所有的请求处理完毕时再退出
	notFoundHandler HandlersChain  // 找不到路由时默认执行的路由
	file            *os.File
	currency int64
}

func (s *Server) stopWgCounter(ctx *Context) {
	s.stopWg.Add(1)
	atomic.AddInt64(&s.currency,1)
	ctx.Next()
	s.stopWg.Done()
	atomic.AddInt64(&s.currency,-1)
}

func (s *Server)Currency()int64{
	return atomic.LoadInt64(&s.currency)
}

func NewServer() *Server {
	s := &Server{
		RouterGroup: RouterGroup{
			path:       "",
			handlers:   nil,
			routerTree: &routerTree{},
		},
	}
	//fmt.Println(logo3)

	s.routerTree.server = s
	//每次接受到一个请求，wg +1 ，处理完一个请求，wg -1
	s.Use(s.stopWgCounter)
	return s
}

func (r *Server) NotFound(handler Handler) {
	r.notFoundHandler = combineHandlers(r.handlers, handler)
}

func (s *Server) request() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		c := newContext()
		c.FastCtx = ctx
		c.Path = stringOf(ctx.Path())
		c.Method = stringOf(ctx.Method())
		c.handlers = nil
		c.RequestURI = stringOf(ctx.RequestURI())
		s.routerTree.handleHTTPRequest(c)
		c.free()
	}
}

func (s *Server) RunFork(addr string, resuport bool) error {
	sv := &fasthttp.Server{
		Handler: s.request(),
	}
	pf := prefork.New(sv)
	pf.Reuseport = resuport
	return pf.ListenAndServe(addr)
}

func getMaxProcs()int{
	num ,_ := strconv.Atoi(os.Getenv("NS_GOMAXPROC"))
	if num > 0{
		return num
	}
	return runtime.NumCPU()
}

func (s *Server) Run(addr string,reusepor bool)(err error ){
	if !reusepor{
		runtime.GOMAXPROCS(getMaxProcs())

	}
	var ls net.Listener
	if reusepor{
		ls ,err  =reuseport.Listen("tcp4",addr)
	}else{
		ls, err = net.Listen("tcp", addr)
	}
	if err != nil {
		return err
	}
	return s.run(ls)
}

func (s *Server) run(ls net.Listener) error {
	// listen

	go func() {
		if err := fasthttp.Serve(ls, s.request()); err != nil {
			if s.closed { // 正常关闭直接return
				return
			} else {
				panic(err) // 否则panic
			}
		}
	}()
	//监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	fmt.Printf("server receive signal:%v ,close listener and start to stop,%v\n", sig, time.Now().String())
	s.closed = true
	// 获取到退出信号时，关闭listener，不再接受新的请求，并且等待所有的请求处理完毕后退出
	ls.Close()
	s.stopWg.Wait()
	fmt.Printf("server successful stoped  %v \n", time.Now().String())
	return nil
}

// 开启多进程
const (
	profkFlags = "-forked-child-process"
	child      = "child"
)

var (
	forkFlagVal = ""
)

func init() {
	flag.StringVar(&forkFlagVal, profkFlags[1:], "", "")
}

func (s *Server) runChild() (cmd *exec.Cmd, err error) {
	cmd = exec.Command(os.Args[0], append(os.Args[1:], profkFlags, child)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{s.file}
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}

func isChild() bool {
	if forkFlagVal == child {
		return true
	}
	return false
}

func (s *Server) listenChild() error {

	ls, err := net.FileListener(os.NewFile(3, ""))
	if err != nil {
		return err
	}
	return s.run(ls)
}

func (s *Server) listenTcpFiles(addr string) error {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	ls, err := net.ListenTCP("tcp", raddr)
	if err != nil {
		return err
	}

	f, err := ls.File()
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

type process struct {
	err error
	pid int
}

func getWorkerProcessesNum()int{
	num ,_ := strconv.Atoi(os.Getenv("NS_WORKER_PROCESSES"))
	if num > 0 && num < 128{
		return num
	}
	return runtime.NumCPU()
}

// 启动子进程进行服务
func (s *Server) forkChildren()error {
	processNum := getWorkerProcessesNum()
	stopChan := make(chan process,processNum)
	childMap := sync.Map{}
	for i := 0; i < processNum; i++ {
		cmd,err := s.runChild()
		if err != nil{
			return err
		}
		go func() {
			stopChan <- process{
				err: cmd.Wait(),
				pid: cmd.Process.Pid,
			}
		}()
		childMap.Store(cmd.Process.Pid,cmd)
	}
	var re  error
	for i:=0 ;i< processNum;i++{
		p := <- stopChan
		log.Println("child process done,err:",p.err,"pid",p.pid)
		childMap.Delete(p.pid)
		cmd,err := s.runChild()
		re = err
		if err != nil{
			return err
		}
		go func() {
			stopChan <- process{
				err: cmd.Wait(),
				pid: cmd.Process.Pid,
			}
		}()
		childMap.Store(cmd.Process.Pid,cmd)
	}
	return re
}

func (s *Server) RunPFork(addr string,reusePort bool) error {
	flag.Parse()
	if isChild() {
		runtime.GOMAXPROCS(1)
		if reusePort{
			ls,err := reuseport.Listen("tcp4",addr)
			if err != nil{
				return err
			}
			return s.run(ls)
		}
		return s.listenChild()
	}
	if !reusePort{
		err := s.listenTcpFiles(addr)
		if err != nil {
			return err
		}
	}

	return  s.forkChildren()
}

type Handler func(ctx *Context)

type HandlersChain []Handler

type RouterGroup struct {
	path       string
	handlers   []Handler
	routerTree *routerTree
}

// 添加拦截器（handlers）
func (r *RouterGroup) Use(handlers ...Handler) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *RouterGroup) Method(method string, pth string, handler Handler) {
	pth = path.Join(r.path, pth)
	r.routerTree.addRoute(method, pth, combineHandlers(r.handlers, handler))
}

//
func (r *RouterGroup) GET(path string, handler Handler) {
	r.Method("GET", path, handler)
}

func (r *RouterGroup) POST(path string, handler Handler) {
	r.Method("POST", path, handler)
}

func (r *RouterGroup) PUT(path string, handler Handler) {
	r.Method("PUT", path, handler)
}

func (r *RouterGroup) DELETE(path string, handler Handler) {
	r.Method("DELETE", path, handler)
}

func (r *RouterGroup) PATCH(path string, handler Handler) {
	r.Method("PATCH", path, handler)
}

func (r *RouterGroup) HEAD(path string, handler Handler) {
	r.Method("HEAD", path, handler)
}

func (r *RouterGroup) OPTION(path string, handler Handler) {
	r.Method("OPTION", path, handler)
}

func (r *RouterGroup) RegisterApis(apis []*Api) *ApiGroup {
	return NewApiGroup(r, apis)
}

func (r *RouterGroup) Any(path string, handler Handler) {
	r.GET(path, handler)
	r.POST(path, handler)
	r.PATCH(path, handler)
	r.DELETE(path, handler)
	r.PUT(path, handler)
	r.HEAD(path, handler)
	r.OPTION(path, handler)
}

//创建一个group，可以构建新的拦截器链路，新的group 会继承父拦截器（handler）
func (r *RouterGroup) Group(pth string) *RouterGroup {
	g := &RouterGroup{
		path:       path.Join(r.path, pth),
		handlers:   combineHandlers(r.handlers), // 这个地方需要复制一份，否则会出现不同的group handler 相互覆盖的情况
		routerTree: r.routerTree,
	}
	return g
}

//recover handler
// 兜底，防止业务逻辑出现panic导致服务直接不可用
func DefaultRecover(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			fmt.Fprintf(os.Stdout, "panic: err:%v stack:%s", err, stringOf(stack))
			server500(c)
			return
		}
	}()
	c.Next()
}

func combineHandlers(hs HandlersChain, handler ...Handler) HandlersChain {
	targets := make(HandlersChain, len(hs), len(hs)+len(handler))
	copy(targets, hs)
	return append(targets, handler...)
}
