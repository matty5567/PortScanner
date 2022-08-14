package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type PortResult struct {
	success bool
	port    string
}

var IPV4_ADDRESS string = "51.183.47.213"
var MAX_PORT int = 64000
var TIMEOUT time.Duration = time.Second
var NUM_WORKERS int = 100

func main() {
	ports_chan := make(chan string, MAX_PORT)
	results_chan := make(chan PortResult)

	// send ports to port_chan
	for port := 1; port <= MAX_PORT; port++ {
		str_port := strconv.Itoa(port)
		ports_chan <- str_port
	}

	// we've added all ports so can close it now
	close(ports_chan)

	// start workers
	for w := 0; w < NUM_WORKERS; w++ {
		go worker(ports_chan, results_chan)
	}

	for i := 1; i < MAX_PORT; i++ {
		portResult := <-results_chan

		if portResult.success == true {
			fmt.Println("Port open on: ", portResult.port)
		}
	}
}

func worker(ports <-chan string, results chan<- PortResult) {
	for port := range ports {

		address := net.JoinHostPort(IPV4_ADDRESS, port)
		conn, _ := net.DialTimeout("tcp", address, TIMEOUT)

		var success bool
		if conn != nil {
			defer conn.Close()
			success = true
		} else {
			success = false
		}
		results <- PortResult{success, port}
	}
}
