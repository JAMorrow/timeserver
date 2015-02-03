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
	"os"
	"net/http"
	"timeserverhtml"
)

const (
	versionNumber = "1.3" // current version number of the software
)

var (
	port = flag.String("port", "8080", "the port number used for the webserver")
	version = flag.Bool("V", false, "display the version number to console")
	templates = flag.String("-templates", "../timeserverhtml/templates/",
		"the directory where the page templates are located.")
 
)

func main() {
        flag.Parse() // get command line arguments

	// check if version number is requested.
	if (*version) {
		fmt.Printf("timeserver Version %s\n", versionNumber)
	}
	
	// Set the templates directory
	timeserverhtml.SetTemplatesDirectory(*templates)

	// Setup handlers for the pages.
	http.HandleFunc("/time", timeserverhtml.TimeHandler)
	http.HandleFunc("/login", timeserverhtml.LoginHandler)
	http.HandleFunc("/logout", timeserverhtml.LogoutHandler)
	http.HandleFunc("/index", timeserverhtml.IndexHandler)
//	http.HandleFunc("/", timeserverhtml.Page404Handler)

	// listen at the given port
	err := http.ListenAndServe(":" + *port, nil)

	// check if there was a problem listening at that port.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
