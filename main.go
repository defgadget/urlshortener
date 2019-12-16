package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strings"
	"regexp"
	"os"
	"encoding/csv"
)
// TODO: Improve UI of website
// TODO: Automate shortening of long url. ?? How to do it?
// TODO: Make code more modular - Move Modal into it's own package
// TODO: Upgrade CSV to a more robust DB

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) save() error {
	
// }

var validPath = regexp.MustCompile("^/([-a-zA-z0-9]+)$")

func expandShortPath(p string) string {
	f, err := os.OpenFile("paths.csv", os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("There was an error opening paths.csv", err)
	}
	r := csv.NewReader(f)
	for {
		rec, err := r.Read() 
		if err != nil {
			fmt.Println("No path found", err)
			return ""
		}
		if rec[0] == strings.ToLower(p) {
			return rec[1]
		}
	}
}
func createShortPath(fp string, sp string) {
	if expandShortPath(sp) != "" {
		fmt.Println("Shortpath already exists")
		return
	}
	f, err := os.OpenFile("paths.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("There was an error opening the file", err)
		return
	}
	w := csv.NewWriter(f)
	err = w.Write([]string{strings.ToLower(sp), fp})
	if err != nil {
		fmt.Println("Error writing to file", err)
		return
	}
	w.Flush()
	fmt.Printf("CREATED: %v will reroute to %v\n", sp, fp)
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