// Thank you Professor Bernstein!

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"github.com/JKowalsky/counter"
	log "github.com/cihub/seelog"
)

var (
	rate = flag.Int("-rate", 200, "average rate of requests (per second)")
	burst = flag.Int("-burst", 30, "number of concurrent requests to issue")
	timeoutMs = flag.Int("-timeout-ms", 400, "max time to wait for response.")
	f_runtime = flag.Int("-rate", 20, " number of seconds to process")
	url = flag.String("--url", "localhost:8080/time", "url to sample.")

	runtime = time.Duration((*f_runtime)) * time.Second // get runtime in seconds
)

var (
	c = counter.New()
)


// A map of replies to count of those replies.
var convert = map[int]string {
	1: "100s",
	2: "200s",
	3: "300s",
	4: "400s",
	5: "500s",
}

// Generate a new request to the server.
func request() {
	log.Info("New Request.")
	timeout := time.Duration((*timeoutMs)) * time.Millisecond
	client := http.Client{
		Timeout : timeout,
	}
	response, err := client.Get(url)
	if err != nil {
		log.Error("No response.")
		c.Incr("total", 1)
		return
	}
	key, ok := convert[response.StatusCode / 100]
	log.Info("Response: ", key)
	if !ok {
		log.Error("Response was an error.")
		key = "errors"
	}
	c.Incr(key, 1)
}

// load and start firing off bursts of requests.
func load() {
	timeout := time.Tick(runtime)
	interval := time.Duration((1000000 * (*burst)) / (*rate)) * time.Microsecond
	period := time.Tick(interval)

	log.Info("Loading.")
	for {
		log.Info("Fire a burst.")
		// fire off burst
		for i := 0; i < (*burst); i++ {
			go request()
		}
		//wait for next tick
		<- period

		// poll for timeout
		select {
		case <-timeout :
			return
		default:

		}	
	}
}

func main () {
	load()
	time.Sleep( time.Duration( 2 * (*timeoutMs)) * time.Millisecond )
	fmt.Printf("total: \t%d\n", c.Get("total"))
	fmt.Printf("100s total: \t%d\n", c.Get("100s"))
	fmt.Printf("200s total: \t%d\n", c.Get("200s"))
	fmt.Printf("300s total: \t%d\n", c.Get("300s"))
	fmt.Printf("400s total: \t%d\n", c.Get("400s"))
	fmt.Printf("500s total: \t%d\n", c.Get("500s"))
	fmt.Printf("total: \t%d\n", c.Get("total"))


}
