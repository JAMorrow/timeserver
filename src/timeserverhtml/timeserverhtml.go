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
	"fmt"
	"net/http"
	"time"
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
func TimeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "<html>")
	fmt.Fprintln(rw, "<head>")
	fmt.Fprintln(rw, "<style>")
	fmt.Fprintln(rw, "p {font-size: xx-large}")
	fmt.Fprintln(rw, "span.time {color: red}")
	fmt.Fprintln(rw, "</style>")
	fmt.Fprintln(rw, "</head>")
	fmt.Fprintln(rw, "<body>")
	fmt.Fprintln(rw, "<p>The time is now <span class=\"time\">")
	fmt.Fprintln(rw, getCurrentTime())
	fmt.Fprintln(rw, "</span>.</p>")
	fmt.Fprintln(rw, "</body>")
	fmt.Fprintln(rw, "</html>")
}

// serves a 404 webpage if the url requested is not found.
func Page404Handler(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
	fmt.Fprintln(rw, "<html>")
	fmt.Fprintln(rw, "<body>")
	fmt.Fprintln(rw, "<p>These are not the URLs you're looking for.</p>")
	fmt.Fprintln(rw, "</body>")
	fmt.Fprintln(rw, "</html>")
}

// serves an index webpage if the user has already logged in.
func IndexHandler(rw http.ResponseWriter, r *http.Request) {

	// check if cookie is set
	cookie, _ := r.Cookie("name")

	// if not, goto login page

	// else say hi

	fmt.Fprintln(rw, "<html>")
	fmt.Fprintln(rw, "<body>")
	fmt.Fprintln(rw, "Greetings, ")
	// TODO name here
	fmt.Fprint(rw, cookie)
	fmt.Fprintln(rw, "</p>")
	fmt.Fprintln(rw, "</body>")
	fmt.Fprintln(rw, "</html>")
}

// serves a Login webpage if the user has not logged in.
func LoginHandler(rw http.ResponseWriter, request *http.Request) {

	fmt.Println("Accessed login")

	username := request.FormValue("name")
	fmt.Println("username is " + username)

	// if name is valid
	if username != "" {
		// set the cookie with the name
		cookie := http.Cookie{Name: "name", Value:username, Path:"/", Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: false}

		http.SetCookie(rw, &cookie)
		//request.AddCookie(&cookie)
		http.Redirect(rw, request, "/index", http.StatusAccepted)

	} else {

	fmt.Fprintln(rw, "<html>")
	fmt.Fprintln(rw, "<body>")
	fmt.Fprintln(rw, "<form action=\"login\">")
	fmt.Fprintln(rw, "What is your name, Earthling?")
	fmt.Fprintln(rw, "<input type=\"text\" name=\"name\" size=\"50\">")
	fmt.Fprintln(rw, "<input type=\"submit\">")
	fmt.Fprintln(rw, "</form>")
	fmt.Fprintln(rw, "</p>")
	fmt.Fprintln(rw, "</body>")
	fmt.Fprintln(rw, "</html>")

	}
	// set cookie
	// else do nothing

}

// serves a Logout webpage if the user has logged in and now wants to logout.
func LogoutHandler(rw http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(rw, "<html>")
	fmt.Fprintln(rw, "<META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprintln(rw, "<body>")
	fmt.Fprintln(rw, "<p>Good-bye.</p>")
	fmt.Fprintln(rw, "</body>")
	fmt.Fprintln(rw, "</html>")
}
