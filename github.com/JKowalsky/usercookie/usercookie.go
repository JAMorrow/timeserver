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
	"sync"
	"time"
	log "github.com/cihub/seelog"
)

var (
	usersUpdating = &sync.Mutex{} // used to lock the users map when adding users
	users = make(map[string]string)
)

func CreateCookie(rw http.ResponseWriter, username string) bool {
	// get unique key via uuidgen
	cmd := exec.Command("uuidgen", "-r") // create a random uuidgen
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error("Error: Unable to run uuidgen.")
		return false
	}

	id := out.String() // the key
	// id has trailing /n, needs to be removed.
	id = strings.TrimSuffix(id, "\n")

	log.Info("Uuidgen for user %s: %s \n", username, id)

	usersUpdating.Lock() // enter mutex while updating users
	users[id] = username
	usersUpdating.Unlock() // exit mutex

	// set the cookie with the name
	cookie := http.Cookie{Name: "Userhash", Value: id, Path: "/", Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: false}

	http.SetCookie(rw, &cookie)
	return true
}

func LogoutCookie(rw http.ResponseWriter, request *http.Request) {
	cookie, _ := request.Cookie("Userhash")
	cookie.MaxAge = -1 // delete the cookie
	cookie.Expires = time.Now()
	cookie.Value = ""          // set the value to null for safety
	http.SetCookie(rw, cookie) // write this to the cookie
}

// return a username associated with a cookie
// TODO: Should return the error if there is any
func GetUsername(r *http.Request) string {
	cookie, err := r.Cookie("Userhash")
	name := ""
	if err == nil { // there is a cookie, print name
		name = users[cookie.Value]
	}
	return name
}

// check if cookie is set
func CookieExists(r *http.Request) bool {
	_, err := r.Cookie("Userhash")
	if err == nil { // there is a cookie
		return true
	} else {
		return false
	}
}
