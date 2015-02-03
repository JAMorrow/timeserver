// timeserver html
// A collection of html serving functions for timeserver
//
// Based on https://golang.org/doc/articles/wiki/final.go
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright @ January 2015, Jennifer Kowalsky

package timeserverhtml

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	loginVisited bool = false // used to keep track of whether or not the login page is visited

	usersUpdating = &sync.Mutex{} // used to lock the users map when adding users

	users = make(map[string]string)

	TemplateDir     = "../timeserverhtml/templates/" // default tempaltes location
	TemplatePage string
	TemplateTime string
	TemplateIndex string
	TemplateLogin string
	TemplateLogout string
	TemplatePage404 string
)

type TimeContext struct {
	Timebody string
	UTCTime  string
	TimeName string
}

type IndexContext struct {
	Username string
}

type LoginContext struct {
	Prompt string
}

type GenericContext struct {
	Unused string
}

// set the templates directory and update all pages to include that path
func SetTemplatesDirectory(templatedir string) {

	fmt.Printf("set templates dir: %s\n", templatedir)

	TemplateDir = templatedir
	TemplatePage    = TemplateDir + "page.html"
	TemplateTime    = TemplateDir + "time.html"
	TemplateIndex   = TemplateDir + "index.html"
	TemplateLogin   = TemplateDir + "login.html"
	TemplateLogout  = TemplateDir + "logout.html"
	TemplatePage404 = TemplateDir + "page404.html"
}


// returns a template containing a specific page's body inside the general page
// Handles errors internally
func makeTemplate(body string) *template.Template {
	tmpl := template.New("page")
	tmpl, err := tmpl.ParseFiles(TemplatePage, body)
	if err != nil {
		fmt.Printf("parsing template: %s\n", err)
		return nil
	}
	return tmpl
}

// Get the current time and return it as a string.
// Note: Removes date and timezone information.
func getCurrentTime() string {
	// layout shows by example how the reference time should be represented.
	const layout string = "3:04:02PM"
	t := time.Now()
	return t.Format(layout)
}

// serves a webpage that returns the current time.
func TimeHandler(rw http.ResponseWriter, r *http.Request) {

	tmpl := makeTemplate(TemplateTime)

	// start building body string
	timeBody := getCurrentTime()

	const layout string = "3:04:02 UTC"
	utcTime := time.Now().UTC().Format(layout)

	// check if cookie is set
	cookie, err := r.Cookie("Userhash")
	timeName := ""
	if err == nil { // there is a cookie, print name
		timeName += ", "
		timeName += users[cookie.Value]
	}

	timecontext := TimeContext{
		Timebody: timeBody,
		UTCTime:  utcTime,
		TimeName: timeName,
	}

	err = tmpl.ExecuteTemplate(rw, "TimeTemplate", timecontext)
	if err != nil {
		fmt.Printf("executing template: %s\n", err)
		return
	}
	fmt.Println("Accessed /time")

}

// serves a 404 webpage if the url requested is not found.
func Page404Handler(rw http.ResponseWriter, r *http.Request) {
	//fmt.Println("Accessed illegal page")
	tmpl := makeTemplate(TemplatePage404)

	context := GenericContext{
		Unused: "",
	}

	err := tmpl.ExecuteTemplate(rw, "Page404Template", context)
	if err != nil {
		fmt.Printf("executing template: %s\n", err)
		return
	}
	http.NotFound(rw, r)
}

// serves an index webpage if the user has already logged in.
func IndexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Accessed /index")
	tmpl := makeTemplate(TemplateIndex)

	// check if cookie is set
	cookie, err := r.Cookie("Userhash")

	if err != nil { // there is no cookie
		http.Redirect(rw, r, "/login", http.StatusBadRequest)

	} else { // else say hi

		indexcontext := IndexContext{
			Username: users[cookie.Value],
		}

		err := tmpl.ExecuteTemplate(rw, "IndexTemplate", indexcontext)
		if err != nil {
			fmt.Printf("executing template: %s\n", err)
			return
		}

	}

}

// serves a Login webpage if the user has not logged in.
func LoginHandler(rw http.ResponseWriter, request *http.Request) {
	fmt.Println("Accessed /login")

	username := request.FormValue("name")
	fmt.Println("username is \"" + username + "\"")

	// sanitize username
	html.EscapeString(username)

	prompt := "" // default is empty

	// if name is valid
	if username != "" && loginVisited {

		// get unique key via uuidgen
		cmd := exec.Command("uuidgen", "-r") // create a random uuidgen
		cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error: Unable to run uuidgen.")
			loginVisited = false
			http.Redirect(rw, request, "/index", http.StatusAccepted)
			return
		}

		id := out.String() // the key
		// id has trailing /n, needs to be removed.
		id = strings.TrimSuffix(id, "\n")

		fmt.Printf("Uuidgen for user %s: %s \n", username, id)

		usersUpdating.Lock() // enter mutex while updating users
		users[id] = username
		usersUpdating.Unlock() // exit mutex

		// set the cookie with the name
		cookie := http.Cookie{Name: "Userhash", Value: id, Path: "/", Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: false}

		http.SetCookie(rw, &cookie)
		loginVisited = false
		http.Redirect(rw, request, "/index", http.StatusAccepted)

	} else if username == "" && loginVisited { // if name is not valid
		prompt = "C'mon, I need a name."
	} else { // first time we hit the page
		loginVisited = true
	}
	tmpl := makeTemplate(TemplateLogin)

	logincontext := LoginContext{
		Prompt: prompt,
	}
	err := tmpl.ExecuteTemplate(rw, "LoginTemplate", logincontext)
	if err != nil {
		fmt.Printf("executing template: %s\n", err)
		return
	}

}

// serves a Logout webpage if the user has logged in and now wants to logout.
func LogoutHandler(rw http.ResponseWriter, request *http.Request) {
	fmt.Println("Accessed /logout")
	// find cookie
	cookie, err := request.Cookie("Userhash")

	if err != nil { // there is no cookie
		http.Redirect(rw, request, "/index", http.StatusBadRequest)

	} else {
		cookie.MaxAge = -1 // delete the cookie
		cookie.Expires = time.Now()
		cookie.Value = ""          // set the value to null for safety
		http.SetCookie(rw, cookie) // write this to the cookie

		tmpl := makeTemplate(TemplateLogout)

		context := GenericContext{
			Unused: "",
		}

		err = tmpl.ExecuteTemplate(rw, "LogoutTemplate", context)
		if err != nil {
			fmt.Printf("executing template: %s\n", err)
			return
		}
	}
}
