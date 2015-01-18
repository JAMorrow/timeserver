// timeserver html
// A collection of html serving functions for timeserver
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
func ServeTime(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<head>")
	fmt.Fprintln(response, "<style>")
	fmt.Fprintln(response, "p {font-size: xx-large}")
	fmt.Fprintln(response, "span.time {color: red}")
	fmt.Fprintln(response, "</style>")
	fmt.Fprintln(response, "</head>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>The time is now <span class=\"time\">")
	fmt.Fprintln(response, getCurrentTime())
	fmt.Fprintln(response, "</span>.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// serves a 404 webpage if the url requested is not found.
func Serve404(response http.ResponseWriter, request *http.Request) {
	http.NotFound(response, request)
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>These are not the URLs you're looking for.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// serves an index webpage if the user has already logged in.
func ServeIndex(response http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "Greetings, ")
	// TODO name here
	fmt.Fprintln(response, "</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// serves a Login webpage if the user has not logged in.
func ServeLogin(response http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<form action=\"login\">")
	fmt.Fprintln(response, "What is your name, Earthling?")
	fmt.Fprintln(response, "<input type=\"text\" name=\"name\" size=\"50\">")
	fmt.Fprintln(response, "<input type=\"submit\">")
	fmt.Fprintln(response, "</form>")
	fmt.Fprintln(response, "</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// serves a Logout webpage if the user has logged in and now wants to logout.
func ServeLogout(response http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>Good-bye.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}
