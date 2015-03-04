// Templates
//
// Manages the templates for the time server
// Uses a single function call to load up a web page and return the template
// to the calling function

package Page

import (
	"html/template"
	"net/http"
	"errors"
)

// Context structs
type LoginContext struct{
	Alert string
}

type HelloContext struct{
	Name string
}

type TimeContext struct{
	Time, Name string
}


var (
	
	directory = "templates" 
	
	// File lists for the different pages
	LoginFiles = []string{
		"/page.tmpl",
		"/login.tmpl",
	}

	HelloFiles = []string{
		"/page.tmpl",
		"/hello.tmpl",
	}

	TimeFiles = []string{
		"/page.tmpl",
		"/time.tmpl",
	}

	YotsubaFiles = []string{
		"/page.tmpl",
		"/yotsuba.tmpl",
	}

	LogoutFiles = []string{
		"/page.tmpl",
		"/logout.tmpl",
	}
	// master file list
	files = map[string][]string{
		"login" : LoginFiles,
		"hello" : HelloFiles,
		"yotsuba" : YotsubaFiles,
		"time" : TimeFiles,
		"logout" : LogoutFiles,
	}
)

func GetPage(w http.ResponseWriter, page string, context interface{})(error){

	tmpl := template.New(page)

	// check for the page
	
	val, ok := files[page]
	if ok {
		var err error
		tmpl, err = tmpl.ParseFiles( directory + val[0], directory + val[1] )
		if err != nil {
			
			// template error
			return errors.New("Template error")
		}

	} else {

		// page not found error
		return errors.New("Page type not found")

	}
	//no errors

	tmpl.ExecuteTemplate(w,"page",context)
	
	return nil
}

func SetDirectory(dir string) {
	directory = dir;
}
