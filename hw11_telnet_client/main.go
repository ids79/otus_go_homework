package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("required parameters are not specified")
	}
	/*signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel)

	go func() {
		for {
			s := <-signalChanel
			fmt.Println(s)
		}
	}()*/

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout of the chat")
	flag.Parse()
	address := net.JoinHostPort(os.Args[len(os.Args)-2], os.Args[len(os.Args)-1])
	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := telnet.Connect(); errors.Is(err, ErrFailedConnection) {
		log.Fatal(err.Error())
	}
	defer telnet.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			if err := telnet.Receive(); err != nil {
				break
			}
		}
		wg.Done()
	}()
	go func() {
		for {
			if err := telnet.Send(); err != nil {
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
