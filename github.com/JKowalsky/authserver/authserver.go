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
	"time"
	"encoding/json"
	"io/ioutil"
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
	dumpfile = flag.String("-dumpfile", "dumpfile", "A file containing the users information for backup purposes.")
	checkptInterval = flag.Int("-checkpoint-interval", 10, "Updates the backup file every X seconds.  default is 10.")
 
)

// map of user names to cookies
var (
	usersUpdating = &sync.Mutex{} // used to lock the users map when adding users
	users = make(map[string]string)
)

// while server exists, wait checkptInterval number of seconds, then perform backup
func backupUsers() {
	for {
		time.Sleep(time.Duration((*checkptInterval)) * time.Second)
		log.Info("Backing up users.")
		// check if dumpfile already exists
		oldfile, err := os.Open((*dumpfile))
		if err == nil {	// if so, rename it.
			oldfile.Close()
			os.Rename((*dumpfile), (*dumpfile) + ".bak")
		}

		// copy current dictionary to reduce time we're tying it up

		backupUsers := make(map[string]string)
		usersUpdating.Lock()
		for key, value := range users {
			backupUsers[key] = value
		}
		usersUpdating.Unlock()

		// encode our dictionary into a JSON file
		jsonmap, err := json.Marshal(backupUsers)

		// lets see it
		log.Info("jsonmap: " + string(jsonmap))
		if err != nil {
			log.Error(err)
		}

		var (
			file *os.File
		)

		if file, err = os.Create((*dumpfile)); err != nil {
			return
		}
		file.WriteString(string(jsonmap))
		file.Close()

		if (!verifyFileContents(backupUsers)) {
			log.Error("Backup unsuccessful!")
		} else {
			log.Info("Backup Verified.")
		}
	}
}

func verifyFileContents(usersMap map[string]string) bool {
	// retrieve data from backup.
	usersOnFile, _ := retrieveUsersFromBackup()

	// quick check, see if the maps have the same length.  If not, they are different.
	if len(usersMap) != len(usersOnFile) {
		return false
	}

	for key, value := range usersMap {
		if usersOnFile[key] != value {
			return false
		}
	}
	return true
}

// assumes file exists
func retrieveUsersFromBackup() (map[string]string, error) {

	userMap := make(map[string]string)

	content, err := ioutil.ReadFile((*dumpfile))
	log.Info("content: " + string(content))

	if err != nil {
		return userMap, err
	}

	// unmarshal


	err = json.Unmarshal(content, &userMap)
	// for debugging

	for key, value := range userMap {
		log.Info("key: " + key + " value: " + value)
	}

	return userMap, err
}

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

	// check if there is a backup file.  If so, use it.
	// check if dumpfile already exists
	backup, err := os.Open((*dumpfile))
	if err == nil {	// if so, load it.
		backup.Close()
		log.Info("Loading Backup.")
		users, _ = retrieveUsersFromBackup()
	}

	go backupUsers()

	// handle get requests
	go http.HandleFunc("/get", getNameHandler)

	// handle set requests
	go http.HandleFunc("/set", setCookieHandler)

	// refuse other requests
	go http.HandleFunc("/*", defaultHandler)

	// listen at the given port
	err = http.ListenAndServe(":" + *authport, nil)

	// check if there was a problem listening at that port.
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
