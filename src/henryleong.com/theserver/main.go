package main

import (
	"bufio"
	"fmt"
	"io"
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

	log.SetOutput(os.Stdout)

	quitCh = make(chan bool)

	go startServer(flags.Port)

	// start the plugin if it's not a plugin
	if !flags.Plugin {
		port := flags.Port + 1
		pluginCmd := startPlugin(port)

		defer func() {
			if pluginCmd != nil {
				log.Printf("Killing process id: %d", pluginCmd.Process.Pid)
				err := pluginCmd.Process.Kill()
				log.Printf("Kill err: %+v", err)
			}
		}()
	}

	<-quitCh

}

func startPlugin(port int) *exec.Cmd {
	cmdFile := os.Args[0]
	args := []string{
		fmt.Sprintf("--plugin=%s", "true"),
		fmt.Sprintf("--port=%d", port),
	}

	log.Printf("Executing %s with %+v", cmdFile, args)

	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	cmd := exec.Command(cmdFile, args...)
	cmd.Stdout = stdoutW
	cmd.Stderr = stderrW

	cmd.Start()

	log.Printf("New process id: %d", cmd.Process.Pid)

	go readPipe("stdout", stdoutR)
	go readPipe("stderr", stderrR)

	return cmd
}

func readPipe(pipeName string, reader *io.PipeReader) {
	buf := bufio.NewReader(reader)
	for {
		line, err := buf.ReadBytes('\n')
		if line != nil {
			log.Printf("[plugin %s] %s", pipeName, string(line))
		}

		if err == io.EOF {
			return
		}
	}

}
