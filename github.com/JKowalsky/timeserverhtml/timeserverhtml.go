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
	"html"
	"html/template"
	"net/http"
	"time"
	"github.com/JKowalsky/usercookie"
	log "github.com/cihub/seelog"

)

var (
	loginVisited bool = false // used to keep track of whether or not the login page is visited

	TemplateDir     = "/templates/" // default templates location
	TemplatePage    string
	TemplateTime    string
	TemplateIndex   string
	TemplateLogin   string
	TemplateLogout  string
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

	log.Info("set templates dir: " + templatedir)

	TemplateDir = templatedir
	TemplatePage = TemplateDir + "page.html"
	TemplateTime = TemplateDir + "time.html"
	TemplateIndex = TemplateDir + "index.html"
	TemplateLogin = TemplateDir + "login.html"
	TemplateLogout = TemplateDir + "logout.html"
	TemplatePage404 = TemplateDir + "page404.html"
}

// returns a template containing a specific page's body inside the general page
// Handles errors internally
func makeTemplate(body string) *template.Template {
	tmpl := template.New("page")
	tmpl, err := tmpl.ParseFiles(TemplatePage, body)
	if err != nil {
		log.Error("parsing template: %s\n", err)
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
	log.Info("Accessed /time")
	tmpl := makeTemplate(TemplateTime)

	// start building body string
	timeBody := getCurrentTime()

	const layout string = "3:04:02 UTC"
	utcTime := time.Now().UTC().Format(layout)

	timeName := ""
	
	log.Info("Check for cookie")
	if usercookie.CookieExists(r) {
		timeName += ", "
		timeName += usercookie.GetUsername(r)

	}
	log.Info("get Context")
	timecontext := TimeContext{
		Timebody: timeBody,
		UTCTime:  utcTime,
		TimeName: timeName,
	}

	log.Info("executing template...")
	err := tmpl.ExecuteTemplate(rw, "TimeTemplate", timecontext)
	if err != nil {
		log.Error("executing template: %s\n", err)
		return
	}
	log.Info("Accessed /time")

}

// serves a 404 webpage if the url requested is not found.
func Page404Handler(rw http.ResponseWriter, r *http.Request) {
	//log.Info("Accessed illegal page")
	tmpl := makeTemplate(TemplatePage404)

	context := GenericContext{
		Unused: "",
	}

	err := tmpl.ExecuteTemplate(rw, "Page404Template", context)
	if err != nil {
		log.Error("executing template: %s\n", err)
		return
	}
	http.NotFound(rw, r)
}

// serves an index webpage if the user has already logged in.
func IndexHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info("Accessed /index")
	tmpl := makeTemplate(TemplateIndex)

	// check if cookie is set
	if usercookie.CookieExists(r) {
		indexcontext := IndexContext{
			Username: usercookie.GetUsername(r),
		}

		err := tmpl.ExecuteTemplate(rw, "IndexTemplate", indexcontext)
		if err != nil {
			log.Error("executing template: %s\n", err)
			return
		}

	} else {
		http.Redirect(rw, r, "/login", http.StatusBadRequest)

	}
}

// serves a Login webpage if the user has not logged in.
func LoginHandler(rw http.ResponseWriter, request *http.Request) {
	log.Info("Accessed /login")

	username := request.FormValue("name")
	log.Info("username is \"" + username + "\"")

	// sanitize username
	html.EscapeString(username)

	prompt := "" // default is empty

	// if name is valid
	if username != "" && loginVisited {

		if usercookie.CreateCookie(rw, username) {
			loginVisited = false
			http.Redirect(rw, request, "/index", http.StatusAccepted)
			return
		} else { // cookie creation was unsuccessful
			prompt = "Cookie not created.  Try again."
		}

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
		log.Error("executing template: %s\n", err)
		return
	}

}

// serves a Logout webpage if the user has logged in and now wants to logout.
func LogoutHandler(rw http.ResponseWriter, request *http.Request) {
	log.Info("Accessed /logout")
	// find cookie
	//cookie, err := request.Cookie("Userhash")

	if usercookie.CookieExists(request) {

		usercookie.LogoutCookie(rw, request)
		tmpl := makeTemplate(TemplateLogout)

		context := GenericContext{
			Unused: "",
		}

		err := tmpl.ExecuteTemplate(rw, "LogoutTemplate", context)
		if err != nil {
			log.Error("executing template: %s\n", err)
			return
		}
	} else { // there is no cookie
		http.Redirect(rw, request, "/index", http.StatusBadRequest)
	}
}