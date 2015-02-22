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
	"fmt"
	"flag"
	"os"
	"time"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"
	"math/rand"
	"sync"
	"github.com/JKowalsky/usercookie"
	"github.com/Curtalius/Page"
	log "github.com/cihub/seelog"
)

const (
	versionNumber = "1.3" // current version number of the software
)

// command line flags
var (
	port = flag.String("port", "8080", "the port number used for the webserver")
	authport = flag.String("authport", "9080", "the port number used for the authorization webserver")
	authhost = flag.String("authhost", "localhost", "hostname of the authorization webserver")
	version = flag.Bool("V", false, "display the version number to console")
	templates = flag.String("-templates", "templates",
		"the directory where the page templates are located.")
	logname = flag.String("-log", "seelog.xml", "the location/name of the log config file.")
	avgRespMs = flag.Int("-avg-response-ms", 100, "the number of milliseconds on average to get the time.")
	deviationMs = flag.Int("-deviation-ms", 10, "the standard deviation from average milliseconds to get time.  Also in milliseconds.") 
  maxInflight = flag.Int("-max-inflight", 0, "the maximum number of requests that can be serviced.  Default is as many as the server can handle without collapsing.")
)

var (
	authserver string
)

var (
	updatingInflight = &sync.Mutex{} // used to lock the users map when adding users
	inflight = 0
)

// add one to inflight
// check if ma-inflight has been reached, if so, return false
// and do not increment
func addInflight() bool {
  updatingInflight.Lock()

	// if maxInflight == 0, don't bother
	if (*maxInflight) == 0 {
		updatingInflight.Unlock()
		return true
  }
	if (inflight == (*maxInflight)) { // refuse
		updatingInflight.Unlock()
		return false
	} else { // update
		inflight++
		updatingInflight.Unlock()
		return true
	}
}

// subtract one from inflight
func subInflight() {
  updatingInflight.Lock()
	inflight--
  updatingInflight.Unlock()
}

// generates a random normally distributed number from a given
// mean and standard deviation.
func randomNormalDelay(mean_ms int, dev_ms int) time.Duration {
	// get random number
	random_ms := rand.NormFloat64() * float64(dev_ms) + float64(mean_ms)
	log.Info("Random wait time: ", random_ms)

	// create time.Duration
	delay := time.Duration(random_ms) * time.Millisecond
	return delay
}

// get a given username associated with a cookie
func getUsername(id string) string {
	name := ""
	// lookup name associated with the cookie
	resp, rerr := http.Get(authserver + "/get?cookie=" + id )
	if (rerr != nil) {
		fmt.Println(rerr)
		return name
	}
	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		fmt.Println(ioerr)
		return name
	}
	log.Info("Name for this cookie: " + string(body) )
	name = string(body)
	resp.Body.Close()
	return name
}


func helloPage(w http.ResponseWriter, r *http.Request) {
  // check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}
	log.Info("HTTP Request: Home page")
	name := ""
	log.Info("Check for cookie")
	if usercookie.CookieExists(r) {

		cookie, _ := r.Cookie("Userhash")
		log.Info("Cookie Userhash: " + cookie.Value)
		name += ", " + getUsername(cookie.Value)
	}

	// write the page
	context := Page.HelloContext{}
	context.Name = name
	err := Page.GetPage(w,"hello",context)
	if err != nil {
		log.Error(err.Error())
	}
  subInflight()
}
// Login Form
func loginPage(w http.ResponseWriter, r *http.Request) {
  // check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}

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
		id, err := usercookie.CreateCookie(w, name)

		if err == nil {
			// created the cookie, now update authserver
			resp, _ := http.PostForm( authserver + "/set", url.Values{"cookie": {id}, "name": {name}} )
			resp.Body.Close()
			http.Redirect(w,r,"/home",http.StatusFound)
			subInflight()
			return

		}	else { // cookie creation was unsuccessful
			log.Error("Cookie not created.  Try again.")
		}
	}
	log.Info("HTTP Request: Login page")
	err := Page.GetPage(w,"login",context)
	if err != nil {
		log.Error(err.Error())
	}
	// login form
	subInflight()
	return
}


// time server
func timeServer(w http.ResponseWriter, r *http.Request) {
	// check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}
	log.Info("HTTP Request: Time Server Page")

	// simulate load by delaying some random amount
	time.Sleep( randomNormalDelay((*avgRespMs), (*deviationMs)))

	// Get the time and format it
	curTime := time.Now()

	// Finish http code
	utcTime := curTime.UTC()

	myTime := curTime.Format("Jan _2 15:04:05") + " (" + utcTime.Format("15:04:05 UTC") + ")"

	// Add name if available
	name := ""
	log.Info("Check for cookie")
	if usercookie.CookieExists(r) {
		cookie, _ := r.Cookie("Userhash")
		name += ", " + getUsername(cookie.Value)

	}
	context := Page.TimeContext{Name:name,Time:myTime}

	err := Page.GetPage(w,"time",context)

	if err != nil {
		log.Error(err.Error())
	}
	subInflight()
}

// 404 error handler function
func yotsuba(w http.ResponseWriter, r *http.Request) {
  // check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}
	w.WriteHeader(http.StatusNotFound)
	//use template
	err := Page.GetPage(w,"yotsuba",nil)
	if err != nil {
		log.Error(err.Error())
	}
	log.Warn("HTTP Status 404: Page not found redirect")
	subInflight()
}
func logout(w http.ResponseWriter, r *http.Request) {
  // check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}
	if usercookie.CookieExists(r) {
		usercookie.LogoutCookie(w, r)
	}
	err := Page.GetPage(w,"logout",nil)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("HTTP Request: Logout Page")
	subInflight()
}

func bannerHandler(w http.ResponseWriter, r *http.Request) {
  // check if we can handle this request
	acceptedRequest := addInflight()
  
	// if refused, return a server error.
	if !acceptedRequest {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("HTTP Status Internal Server Error")
	}
	file, err := os.Open("bin/banner.css")
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()
	// Copy( writer, reader )
	io.Copy(w,file)

	log.Info("HTTP Request: Banner Style Sheet")
	subInflight()
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
	
	authserver = "http://" + (*authhost) + ":" + (*authport)

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
		fmt.Println("Port occupied!")
		os.Exit(1)
	}
}