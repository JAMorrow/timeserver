package main

import (
	"fmt"
	"html/template"
	"os"
)

var (
	LoginFiles = []string {
		"templates/page.tmpl",
		"templates/login.tmpl",
	}
	HelloFiles = []string{
		"templates/page.tmpl"
		"templates/hello.tmpl"
	}
)
type Context struct {
	Alert string
}

func main() {

	tmpl := template.New("login")
	fmt.Println(tmpl.Name())

/*	tmpl, _ = tmpl.ParseFiles("templates/login.tmpl")
	
	tmpl2 := tmpl.New("login")
	fmt.Println(tmpl2.Name())

//	tmpl3, _ := tmpl2.ParseFiles("templates/page.tmpl")
*/	tmpl3, _ := tmpl.ParseFiles(LoginFiles...)
	fmt.Println(tmpl3.Name())
		
	context1 := Context{Alert:"blah"}
	_ = tmpl3.ExecuteTemplate(os.Stdout, "page", context1)

	tmpl4 := tmpl.Lookup("login")
	if tmpl4 == nil {
		fmt.Println("not found")
	} else {
		fmt.Println("found")
	}

}
