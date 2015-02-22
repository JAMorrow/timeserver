// Authserver
// A server that maintains a database of users.  Associated with timeserver.
//
// Command line arguments: 
// -V displays the version number in the console;
// --port <PORTNUMBER> binds the server to the specified port. 
// 9080 is the default if no port number is given.
//
// Copyright @ January 2015, Jennifer Kowalsky

package main

import (
	"flag"
	"os"
	"net/http"
	"sync"
//	"github.com/JKowalsky/usercookie"
	log "github.com/cihub/seelog"
)

const (
	versionNumber = "0.1" // current version number of the software
)


// command line arguments
var (
	authport = flag.String("port", "9080", "the port number used for the webserver")
	version = flag.Bool("V", false, "display the version number to console")
	templates = flag.String("-templates", "templates",
		"the directory where the page templates are located.")
	logname = flag.String("-log", "seelog.xml", "the location/name of the log config file.")
 
)

// map of user names to cookies
var (
	usersUpdating = &sync.Mutex{} // used to lock the users map when adding users
	users = make(map[string]string)
)

func getNameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("cookie") // the uuid
	log.Info("Cookie id: " + id)
	//name, err := usercookie.GetUsername(r)

	// search our users for this uuid
	if name, ok := users[id]; ok  {
		log.Info("Name associated with this cookie: " + name)
		// convert name to bytes
		w.Write([]byte(name))
		w.WriteHeader(http.StatusOK)

		//w.Header().Add("Name", name)
	} else {
		log.Info("Name not found.")
		w.WriteHeader(http.StatusBadRequest)
	}

	/*if err == nil {
		w.WriteHeader(http.StatusOK)
		// convert name to bytes
		w.Write([]byte(name))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}*/

}


func setCookieHandler (w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("cookie") // the uuid
	name := r.FormValue("name")    // the username

	log.Info("Cookie id: " + id)
	log.Info("Cookie Name: " + name)
	// add pair to our map
	usersUpdating.Lock() // enter mutex while updating users
	users[id] = name
	usersUpdating.Unlock() // exit mutex

	// check that name was added successfully
	if _, ok := users[id]; ok  { 
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func defaultHandler (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	log.Warn("HTTP Status 404: Page not found redirect")

}

func main() {

	flag.Parse() // get command line arguments

	// check if version number is requested.
	if (*version) {
		log.Info("timeserver Version %s\n", versionNumber)
	}

	// handle get requests
	go http.HandleFunc("/get", getNameHandler)

	// handle set requests
	// todo, should dynamically match  /set?cooke=cookie&name=name
	go http.HandleFunc("/set", setCookieHandler)

	// refuse other requests
	go http.HandleFunc("/*", defaultHandler)

	// listen at the given port
	err := http.ListenAndServe(":" + *authport, nil)

	// check if there was a problem listening at that port.
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
