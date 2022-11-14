package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Shortening struct {
	Shortener []Shortener
}

type Shortener struct {
	Key  string `json:"key"`
	Dest []byte `json:"dest"`
}

func (s *Shortening) FindShortener(key string) {
	fmt.Println(s)
}

func (s *Shortening) Exe() {
	fmt.Println(s)

	return
}

func (s *Shortening) Add() {

}

func (s *Shortening) Edit() {

}

func (s *Shortening) Save() {

}

func (s *Shortening) Help() {

}

var shortenings Shortening

func init() {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("dbfile.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Do log if error failed open / create dbfile if not exist
	if err != nil {
		log.Fatal(err)
	}

	jsonValue, _ := ioutil.ReadAll(f)
	fmt.Println(len(jsonValue))

	switch {
	case len(jsonValue) > 0:
		{
			// we unmarshal our byteArray which contains our
			// jsonFile's content into 'Shortening' which we defined above
			errDecod := json.Unmarshal(jsonValue, shortenings)
			if errDecod != nil {
				log.Fatal(errDecod)
			}
		}
	default:
		{
			shortenings = Shortening{}
		}
	}

	return
}

var validShortenerPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")
var validEditorPath = regexp.MustCompile("^/(add|edit|remove)/([a-zA-Z0-9]+)$")

func wrapperEditorHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := validEditorPath.FindStringSubmatch(r.URL.Path)
		if url == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, url[2])
	}
}

func shorteningHandler(w http.ResponseWriter, r *http.Request) {

	uri := validShortenerPath.FindStringSubmatch(r.URL.Path)
	if uri == nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Hi there, Shortening %s", uri[1])

	//
	shortenings.FindShortener(uri[1])

}

func addHandler(w http.ResponseWriter, r *http.Request, key string) {

	//

	p := &Shortener{Key: key}

	html, err := template.ParseFiles("shortenerForm.html")
	if err != nil {
		log.Fatal(err)
	}

	html.Execute(w, p)

}

func editHandler(w http.ResponseWriter, r *http.Request, key string) {

	//

	p := &Shortener{Key: key}

	html, err := template.ParseFiles("shortenerForm.html")
	if err != nil {
		log.Fatal(err)
	}

	html.Execute(w, p)

}

func removeHandler(w http.ResponseWriter, r *http.Request, key string) {

	fmt.Fprintf(w, "Hi there, remove Shortening %s", key)

}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, save Shortening")
}

func helpHandler(w http.ResponseWriter, r *http.Request) {

	//
	fmt.Fprintf(w, "%-40s%s", "/[YOUR-SHORTENER-KEY]", "URI that your shortener works. \n")
	fmt.Fprintf(w, "%-40s%s", "", "If your shortener exist, you will redirected to destination location. \n")
	fmt.Fprintf(w, "%-40s%s", "", "And if not exist, you will be redirected to add shortener.\n\n")

	//
	fmt.Fprintf(w, "%-40s%s", "/add/[YOUR-SHORTENER-KEY]", "URI that you can add new shortener. \n")
	fmt.Fprintf(w, "%-40s%s", "", "If your shortener not exist, you can add new shortener. \n")
	fmt.Fprintf(w, "%-40s%s", "", "if exist, you will be redirected to edit shortener.\n\n")

	//
	fmt.Fprintf(w, "%-40s%s", "/edit/[YOUR-SHORTENER-KEY]", "URI that you can edit shortener. \n")
	fmt.Fprintf(w, "%-40s%s", "", "If your shortener exist, you can edit shortener. \n")
	fmt.Fprintf(w, "%-40s%s", "", "if not exist, you will be redirected to add shortener.\n\n")

	//
	fmt.Fprintf(w, "%-40s%s", "/remove/[YOUR-SHORTENER-KEY]", "URI that you can remove shortener. \n\n")

	//
	fmt.Fprintf(w, "%-40s%s", "/help", "Your current URI that show you URI information. \n\n")

}

func main() {

	//
	http.HandleFunc("/", shorteningHandler)

	//
	http.HandleFunc("/save", saveHandler)

	//
	http.HandleFunc("/add/", wrapperEditorHandler(addHandler))

	//
	http.HandleFunc("/edit/", wrapperEditorHandler(editHandler))

	//
	http.HandleFunc("/remove/", wrapperEditorHandler(removeHandler))

	//
	http.HandleFunc("/help", helpHandler)

	//
	log.Fatal(http.ListenAndServe(":8080", nil))

}
