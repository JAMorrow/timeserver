// Timeserver
// A server that generates a webpage displaying the current time.
//
// Command line arguments: 
// -V displays the version number in the console;
// --port <PORTNUMBER> binds the server to the specified port. 
// 8080 is the default if no port number is given.
//
// Copyright @ January 2015, Jennifer Kowalsky

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"os"
)

const (
	versionNumber = "1.0" // current version number of the software
)

var (
	port = flag.String("port", "8080", "the port number used for the webserver")
	version = flag.Bool("V", false, "display the version number to console")
 
)

// Get the current time and return it as a string.
// Note: Removes date and timezone information.
func getCurrentTime() string {
        // layout shows by example how the reference time should be represented.
	const layout string = "3:04:02PM"
        t := time.Now()
        return t.Format(layout)
}

// serves a webpage that returns the current time.
func serveTime(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<head>")
	fmt.Fprintln(response, "<style>")
	fmt.Fprintln(response, "p {font-size: xx-large}")
	fmt.Fprintln(response, "span.time {color: red}")
	fmt.Fprintln(response, "</style>")
	fmt.Fprintln(response, "</head>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>The time is now <span class=\"time\">")
	fmt.Fprintln(response, getCurrentTime())
	fmt.Fprintln(response, "</span>.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// serves a 404 webpage if the url requested is not found.
func serve404(response http.ResponseWriter, request *http.Request) {
	http.NotFound(response, request)
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>These are not the URLs you're looking for.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

func main() {
        flag.Parse() // get command line arguments

	// check if version number is requested.
	if (*version) {
		fmt.Printf("timeserver Version %s\n", versionNumber)
	}

	// Setup handlers for the pages.
	http.HandleFunc("/time", serveTime)
	http.HandleFunc("/", serve404)

	// listen at the given port
	err := http.ListenAndServe(":" + *port, nil)

	// check if there was a problem listening at that port.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
