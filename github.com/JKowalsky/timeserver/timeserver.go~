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
	"os"
	"time"
	"io"
	"net/http"
	"strings"
	"github.com/JKowalsky/usercookie"
	"github.com/Curtalius/Page"
	log "github.com/cihub/seelog"
)

const (
	versionNumber = "1.3" // current version number of the software
)

var (
	port = flag.String("port", "8080", "the port number used for the webserver")
	version = flag.Bool("V", false, "display the version number to console")
	templates = flag.String("-templates", "templates",
		"the directory where the page templates are located.")
	logname = flag.String("-log", "seelog.xml", "the location/name of the log config file.")
 
)


func helloPage(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP Request: Home page")
	name := ""
	log.Info("Check for cookie")
	if usercookie.CookieExists(r) {
		name += ", "
		name += usercookie.GetUsername(r)

	}

	// write the page
	context := Page.HelloContext{}
	context.Name = name
	err := Page.GetPage(w,"hello",context)
	if err != nil {
		log.Error(err.Error())
	}
}
// Login Form
func loginPage(w http.ResponseWriter, r *http.Request) {
	if usercookie.CookieExists(r) {
		http.Redirect(w,r,"/home",http.StatusFound)
	}

	r.ParseForm()
	name := r.FormValue("name")
	submit := r.FormValue("submit")
	
	// Make a context for the login page
	context := Page.LoginContext{}
	if name == "" && submit == "Submit" {
		log.Warn("Blank name given")
		context.Alert = "Cmon I need a name"
	} else if submit == "Submit" && strings.ContainsAny( name, "#&%/" ) {
		log.Warn("Illegal name given")
		context.Alert = "Names can't contain the following characters #&%/"
	} else if name != "" {
		if usercookie.CreateCookie(w, name) {
			http.Redirect(w, r, "/index", http.StatusAccepted)
			return
		} else { // cookie creation was unsuccessful
			log.Error("Cookie not created.  Try again.")
		}
	}
	log.Info("HTTP Request: Login page")
	err := Page.GetPage(w,"login",context)
	if err != nil {
		log.Error(err.Error())
	}
	// login form
	return
}
// time server
func timeServer(w http.ResponseWriter, r *http.Request) {

	log.Info("HTTP Request: Time Server Page")

	// Get the time and format it
	curTime := time.Now()

	// Finish http code
	utcTime := curTime.UTC()

	myTime := curTime.Format("Jan _2 15:04:05") + " (" + utcTime.Format("15:04:05 UTC") + ")"

	// Add name if avail

	name := ""
	log.Info("Check for cookie")
	if usercookie.CookieExists(r) {
		name += ", "
		name += usercookie.GetUsername(r)

	}
	context := Page.TimeContext{Name:name,Time:myTime}

	err := Page.GetPage(w,"time",context)

	if err != nil {
		log.Error(err.Error())
	}
}

// 404 error handler function
func yotsuba(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	//use template
	err := Page.GetPage(w,"yotsuba",nil)
	if err != nil {
		log.Error(err.Error())
	}
	log.Warn("HTTP Status 404: Page not found redirect")
}
func logout(w http.ResponseWriter, r *http.Request) {

	if usercookie.CookieExists(r) {
		usercookie.LogoutCookie(w, r)
	}
	err := Page.GetPage(w,"logout",nil)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("HTTP Request: Logout Page")

}

func bannerHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("bin/banner.css")
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()
	// Copy( writer, reader )
	io.Copy(w,file)

	log.Info("HTTP Request: Banner Style Sheet")
}



func main() {

	/*logger, err := log.LoggerFromConfigAsFile(*logname)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.ReplaceLogger(logger)*/


	defer log.Flush()
      
	flag.Parse() // get command line arguments

	// check if version number is requested.
	if (*version) {
		log.Info("timeserver Version %s\n", versionNumber)
	}
	

	// Setup handlers for the pages.

	Page.SetDirectory(*templates)

	// Style Sheet
	go http.HandleFunc("/banner.css", bannerHandler)

	// Logout page
	go http.HandleFunc("/logout", logout)

	// Time message
	go http.HandleFunc("/time", timeServer)

	// Login Page handler
	go http.HandleFunc("/", loginPage)
	go http.HandleFunc("/login", loginPage)

	// Home Page handler
	go http.HandleFunc("/home", helloPage)
	go http.HandleFunc("/index", helloPage)


	// listen at the given port
	err := http.ListenAndServe(":" + *port, nil)

	// check if there was a problem listening at that port.
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
