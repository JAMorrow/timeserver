// user cookie
// Code to support a user listing for the timeserver.
//
// Based on https://golang.org/doc/articles/wiki/final.go
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright @ Febuary 2015, Jennifer Kowalsky

package usercookie

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
	"time"
	log "github.com/cihub/seelog"
)

// Create a cookie
// return the generated id and whether or not we were successful
func CreateCookie(rw http.ResponseWriter, username string) (string, error) {
	// get unique key via uuidgen
	id := ""
	cmd := exec.Command("uuidgen", "-r") // create a random uuidgen
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error("Error: Unable to run uuidgen.")
		return id, err
	}	
	id = out.String() // the key
	// id has trailing /n, needs to be removed.
	id = strings.TrimSuffix(id, "\n")

	log.Info("Uuidgen for user " + username + ": " + id)

	/*usersUpdating.Lock() // enter mutex while updating users
	users[id] = username
	usersUpdating.Unlock() // exit mutex*/

	// set the cookie with the id
	cookie := http.Cookie{Name: "Userhash", Value: id, Path: "/", Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: false}

	http.SetCookie(rw, &cookie)
	return id, err
}

func LogoutCookie(rw http.ResponseWriter, request *http.Request) {
	cookie, _ := request.Cookie("Userhash")
	cookie.MaxAge = -1 // delete the cookie
	cookie.Expires = time.Now()
	cookie.Value = ""          // set the value to null for safety
	http.SetCookie(rw, cookie) // write this to the cookie
}

// return a username associated with a cookie
/*func GetUsername(r *http.Request) (string, error) {
	cookie, err := r.Cookie("Userhash")
	log.Info("Cookie Userhash: " + cookie.Value)
	name := ""
	if err == nil { // there is a cookie, print name
		name = users[cookie.Value]
	}
	return name, err
}*/

// check if cookie is set
func CookieExists(r *http.Request) bool {
	_, err := r.Cookie("Userhash")
	if err == nil { // there is a cookie
		return true
	} else {
		return false
	}
}
