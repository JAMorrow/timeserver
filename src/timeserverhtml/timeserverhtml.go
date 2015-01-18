// timeserver html
// A collection of html serving functions for timeserver
//
// Copyright @ January 2015, Jennifer Kowalsky

package timeserverhtml

import (
	"time"
	"net/http"
	"fmt"
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
