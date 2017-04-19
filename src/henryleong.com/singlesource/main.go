package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
)

var quitCh chan bool

func init() {
	flag.BoolVar(&flags.Plugin, "plugin", false, "Start as plugin")
	flag.IntVar(&flags.Port, "port", 7000, "HTTP Server port")

	flag.Parse()
}

type Flags struct {
	Plugin bool
	Port   int
}

var flags *Flags = &Flags{}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			buf := make([]byte, 1<<20)
			buf = buf[:runtime.Stack(buf, true)]
			log.Printf("Caught SIG %+v\n%s", sig, string(buf))

			quitCh <- true
			break
		}
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

	cmd := exec.Command(cmdFile, args...)

	stdoutR, stdoutW := io.Pipe()
	cmd.Stdout = stdoutW
	go readPipe("stdout", stdoutR)

	stderrR, stderrW := io.Pipe()
	cmd.Stderr = stderrW
	go readPipe("stderr", stderrR)

	cmd.Start()

	log.Printf("New process id: %d", cmd.Process.Pid)
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

func startServer(port int) {
	if flags.Plugin {

		var counter = 0

		http.HandleFunc("/greetings", func(w http.ResponseWriter, r *http.Request) {
			counter++
			str := fmt.Sprintf("greetings[%d]:%d\n", port, counter)

			log.Printf("%s", str)

			w.Write([]byte(str))
		})

	} else {
		log.Print("setting up reverse proxy")

		pluginPort := flags.Port + 1
		proxyUrl, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", pluginPort))

		reverseProxy := httputil.NewSingleHostReverseProxy(proxyUrl)
		http.HandleFunc("/greetings", reverseProxy.ServeHTTP)
	}

	bind := fmt.Sprintf(":%d", port)

	log.Printf("Binding server at %s", bind)
	http.ListenAndServe(bind, nil)
}
