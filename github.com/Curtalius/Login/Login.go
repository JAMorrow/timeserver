//Login Package
//
// Manages the Logged in users as well as checking the cookies that
// are related to those logins
package Login

import (
	"fmt"
	"sync"
	"errors"
	"net/http"
	"os/exec"
	"bytes"
	"strings"
	"time"
)

// Structure for maintaining Authentication of users
type Users struct{
	Dict    map[string]string
	Lock    sync.RWMutex
}

var (
	users Users
)
// checks the list for the given name, returns a name and nil if found
// returns "" and an error if not found
func (Users) CheckUser(r *http.Request) (string, error) {
	fmt.Println("Entering CheckUser.")
	// get the cookie
	cookie, err := r.Cookie("UID")
	if err == nil {
		// check the map for the user
		users.Lock.RLock()

		if users.Dict == nil {
			users.Dict = make(map[string]string)
		}
		val, ok := users.Dict[cookie.Value]
		users.Lock.RUnlock()
		if ok {
			// user found
			return val, nil
			
		} else {

			// no user found with found cookie
			return "", errors.New("No user found")
		}
	}
	// the function should have returned by, so return error
	return "", errors.New("No cookie found\r\n")

}

// adds the given name to the map, and adds the cookie to the http client
// It is assumed that the input name is valid at this point in the code
func (Users) AddUser(rw http.ResponseWriter, name string) (error) {
	fmt.Println("Entering AddUser.")

	// get unique key via uuidgen
	cmd := exec.Command("uuidgen", "-r") // create a random uuidgen
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		//log.Error("Error: Unable to run uuidgen.")
		return err
	}

	id := out.String() // the key
	// id has trailing /n, needs to be removed.
	id = strings.TrimSuffix(id, "\n")

	//log.Info("Uuidgen for user %s: %s \n", name, id)

/*	users.Lock() // enter mutex while updating users
	users[id] = username
	usersUpdating.Unlock() // exit mutex
*/
	// set the cookie with the name
	cookie := http.Cookie{Name: "Userhash", Value: id, Path: "/", Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: false}

	http.SetCookie(rw, &cookie)

	if err == nil {
		// Set up cookie jar
		if users.Dict == nil {
			users.Dict = make(map[string]string)
		}

		users.Lock.Lock()
		users.Dict[id] = name
		users.Lock.Unlock()

		return nil
	} else {
		return errors.New("Uuidgen generation error")
	}

}

func (Users) ClearCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering ClearCookie.")
	// Determine if a cookie is present
	cookie, err := r.Cookie("UID")
	
	if err == nil {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)

	// clear out name
	}


}
