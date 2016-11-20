package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var quitCh chan bool

func init() {
	initSignalHandler()
}

func initSignalHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quitCh <- true
	}()
}

func main() {

	var pluginCmd *exec.Cmd = nil

	quitCh = make(chan bool)

	go startServer(flags.Port)

	// start the plugin if it's not a plugin
	if !flags.Plugin {
		port := flags.Port + 1
		pluginCmd = startPlugin(port)
	}

	<-quitCh

	if pluginCmd != nil {
		log.Printf("Killing process id: %d", pluginCmd.Process.Pid)
		err := pluginCmd.Process.Kill()
		log.Printf("Kill err: %+v", err)
	}

}

func startPlugin(port int) *exec.Cmd {
	cmdFile := os.Args[0]
	args := []string{
		fmt.Sprintf("--plugin=%s", "true"),
		fmt.Sprintf("--port=%d", port),
	}

	log.Printf("Executing %s with %+v", cmdFile, args)

	cmd := exec.Command(cmdFile, args...)
	cmd.Start()

	log.Printf("New process id: %d", cmd.Process.Pid)

	return cmd
}
