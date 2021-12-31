package pforks

import (
	"flag"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

const (
	forkFlag = "-fork-process-flag"
	child    = "children"
)

var childFlag = flag.String(forkFlag[1:], "", "indicate if process is master or child")

func IsChildren() bool {
	return *childFlag == child
}

type Fork struct {
	DoMaster   func() error
	DoChildren func() error
	// extra files will be sent to child processes
	Files      []*os.File
	ChildNum   int
	cmds       []*exec.Cmd
	// extra flags will set as cmd args to child processes
	ExtraFlags []string
}

func (f *Fork) runCmd() (*exec.Cmd, error) {
	args :=  append(os.Args[1:], forkFlag, child)
	args = append(args,f.ExtraFlags...)
	cmd := exec.Command(os.Args[0],args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = f.Files
	err := cmd.Start()
	return cmd, err
}

func (f *Fork)Cmds()[]*exec.Cmd{
	return f.cmds
}

func (f *Fork) forkChildren() error{
	if f.ChildNum <= 0 {
		f.ChildNum = runtime.GOMAXPROCS(0)
	}
	if err := f.DoMaster();err != nil{
		panic(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < f.ChildNum; i++ {
		cmd ,err := f.runCmd()
		wg.Add(1)
		if err != nil{
			return err
		}
		f.cmds = append(f.cmds,cmd)
		go func() {
			err = cmd.Wait()
			if err != nil{

			}
			wg.Done()
		}()
	}

	wg.Wait()

	return nil

}

func (f *Fork) Run() error {
	flag.Parse()
	if IsChildren(){
		return  f.DoChildren()
	}
	return  f.forkChildren()
}
