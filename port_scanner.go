package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

type PortResult struct {
	success bool
	port    string
}

var maxPort int
var TIMEOUT time.Duration = time.Second
var numWorkers int = 100
var targetIP string

func main() {

	flag.StringVar(&targetIP, "targetIP", "127.0.0.1", "ip address to run scanner against")
	flag.IntVar(&maxPort, "maxPort", 1024, "Maximum port to scan up to")
	flag.IntVar(&numWorkers, "numWorkers", 50, "Number of goroutines to run")
	flag.Parse()

	fmt.Printf("Running port scanner on ip: %s, up to port %d with %d workers\n", targetIP, maxPort, numWorkers)
	ports_chan := make(chan string, maxPort)
	results_chan := make(chan PortResult)

	// send ports to port_chan
	for port := 1; port <= maxPort; port++ {
		str_port := strconv.Itoa(port)
		ports_chan <- str_port
	}

	// we've added all ports so can close it now
	close(ports_chan)

	// start workers
	for w := 0; w < numWorkers; w++ {
		go worker(ports_chan, results_chan)
	}

	for i := 1; i < maxPort; i++ {
		portResult := <-results_chan

		if portResult.success == true {
			fmt.Println("Port open on: ", portResult.port)
		}
	}
}

func worker(ports <-chan string, results chan<- PortResult) {
	for port := range ports {

		address := net.JoinHostPort(targetIP, port)
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
