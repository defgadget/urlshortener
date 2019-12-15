package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"regexp"
)
// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) save() error {
	
// }

var validPath = regexp.MustCompile("^/([-a-zA-z0-9]+)$")
var urls = make(map[string] string)
func expandShortPath(p string) string {
	// TODO: Make sure that check is case-insensitive. 
	fp, ok := urls[p]
	if !ok {
		fmt.Printf("The short path '%v' cannot be expanded\n", p)
		return ""
	}
	return fp
}
func createShortPath(fp string, sp string) {
	if _, ok := urls[sp]; ok {
		fmt.Println("Shortpath already exists")
		return
	}
	fmt.Printf("CREATED: %v will reroute to %v\n", sp, fp)
	urls[sp] = fp
}
// RouteHandler Expects shortened URL, and looks for a full URL to reroute to.
func RouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		fp := r.Form.Get("FullPath")
		sp := r.Form.Get("ShortPath")
		fmt.Printf("Form posted FullPath: %v, ShortPath: %v\n", fp, sp)
		createShortPath(fp, sp)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.URL.Path == "/" {
		templ, err := template.ParseFiles("form.html")
		if err != nil {
			fmt.Fprint(w, "You have reached a place of no return.")
		}
		err = templ.Execute(w, "")
		if err != nil {
			fmt.Fprint(w, "No ability to execute template")
		}
		return
	}
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		fmt.Println("Not a valid path", r.URL.Path)
		return
	}
	
	shortpath := m[1]
	fullpath := expandShortPath(shortpath)
	if fullpath == "" {
		return
	}
	fmt.Printf("%v is being rerouted to %v\n", shortpath, fullpath)
	http.Redirect(w, r, "/"+fullpath, http.StatusFound)
}
func longURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You've reached the long-url")
}
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	
}
func main() {
	createShortPath("long-url", "long")
	http.HandleFunc("/", RouteHandler)
	http.HandleFunc("/favicon.ico", handleFavicon)
	http.HandleFunc("/long-url", longURL)
	fmt.Println("Serving on 127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}