// TODO
// daemonize itself to run a notification in x time
// [x] daemonize itself
// [x] notify
// [x] flags
// [ ] telegram notifications?
package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jessevdk/go-flags"
	naturaldate "github.com/tj/go-naturaldate"
)

type opt struct {
	Until string `short:"t" long:"time" description:"Notifies after a certain time, Ex: 10 minutes from now, 6:00pm, etc."`
}

func main() {
	// fork
	if os.Args[len(os.Args)-1] != "notforked" {
		tmpfile, _ := exec.LookPath(os.Args[0])
		self := syscall.ProcAttr{
			Dir:   "",
			Env:   []string{},
			Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
			Sys:   nil,
		}
		_, err := syscall.ForkExec(tmpfile, append(os.Args, "notforked"), &self)
		if err != nil {
			panic(err)
		}
		return
	}

	// parse CLI flags
	opts := opt{}
	args, err := flags.Parse(&opts)
	args = args[:len(args)-1]
	if err != nil {
		panic(err)
	}
	if len(args) == 0 {
		return
	}
	var until time.Time
	tmr, err := naturaldate.Parse(opts.Until, time.Now(), naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		panic(err)
	}
	until = tmr

	msg := strings.Join(args, " ")
	<-time.NewTimer(time.Until(until)).C
	err = beeep.Notify(msg, msg, "")
	if err != nil {
		panic(err)
	}
}
