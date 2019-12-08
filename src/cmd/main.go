package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// example -ldflags "-X 'main.BuildVersion=v2' -X 'main.BuildTime=20191202' -X 'main.BuildSum=xxxx'"
var (
	// BuildVersion ...
	BuildVersion string
	// BuildTime ...
	BuildTime string
	// BuildSum ...
	BuildSum string
)

var bdaemon = flag.Bool("daemon", false, "daemon")
var configPath = flag.String("config", "../config.yml", "input config")
var stdrewrite = flag.String("stdrewrite", "./daemon.out", "rewrite stdin stdout stderr")

func init() {
	args := os.Args
	if args != nil && len(args) > 1 && args[1] == "-v" {
		fmt.Printf("version:%s time:%s sum:%s \n", BuildVersion, BuildTime, BuildSum)
		os.Exit(0)
	}
}

func rewriteStd() *os.File {
	f, err := os.OpenFile(*stdrewrite, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	os.Stdin = f
	os.Stdout = f
	os.Stderr = f
	return f
}

func daemon() {

	args := make([]string, 0, len(os.Args))
	for _, item := range os.Args {
		if !strings.HasPrefix(item, "-daemon") {
			args = append(args, item)
		}
	}

	cmd := exec.Command(os.Args[0], args[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} linux
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	err := cmd.Start()
	if err == nil {
		cmd.Process.Release()
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	f := rewriteStd()
	defer f.Close()

	if *bdaemon {
		daemon()
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	//signal.Ignore(syscall.SIGCHLD, syscall.SIGPIPE, syscall.SIGHUP)
	signal.Ignore(syscall.SIGPIPE, syscall.SIGHUP)

	server := NewServer()

	if err := server.Start(*configPath); err != nil {
		fmt.Printf("%s, start error[%s]\n", time.Now().String(), err.Error())
		return
	}
	defer server.Stop()

	fmt.Printf("%s, start success\n", time.Now().String())
	defer func() {
		fmt.Printf("%s, stop\n", time.Now().String())
	}()

	for {
		<-c
		// kill -15 ...reload
		// if s == syscall.SIGTERM {
		// 	  continue
		// }
		break

	}

}
