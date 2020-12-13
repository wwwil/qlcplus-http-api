package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

    "github.com/gorilla/mux"
    "github.com/wwwil/qlcplus-websocket-http/pkg/qlcplus"
)

// Version is the version of the app. This is injected during build.
var Version = "development"

// Commit is the commit hash of the build. This is injected during build.
var Commit string

// BuildDate is the date of the build. This is injected during build.
var BuildDate string

// GoVersion is the Go version used for the build. This is injected during build.
var GoVersion string

// Platform is the target platform for this build. This is injected during build.
var Platform string

// Variables set by argument flags.
var qlcplusAddress = flag.String("qlcplus", "localhost:9999", "Address of QLC+ websocket")
var httpAddress = flag.String("http", "localhost:8888", "Address for QLC+ HTTP API")

var homePageHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width,initial-scale=1">
		<style>
			table, th, td {
				border: 1px solid black;
			}
		</style>
        <title>qlcplus-http-api</title>
    </head>
    <body>
		<h1><a href="https://github.com/wwwil/qlcplus-http-api">qlcplus-http-api</a></h1>
		<h2>Widgets</h2>
		<table>
		<tr>
		  <th>ID</th>
		  <th>Name</th>
		</tr>
		%s
		</table>
    </body>
</html>
`

var widgetHTMLTemplate = `
<tr>
  <td>%s</td>
  <td>%s</td>
</tr>
`

// homePage provides a human readable summary of endpoints.
func homePage(w http.ResponseWriter, r *http.Request) {
	widgetsMap, err := q.GetWidgetsMap()
	if err != nil {
		log.Println(err)
	}
	widgetsHTML := ""
	for widgetID, widgetName := range widgetsMap {
		widgetsHTML = widgetsHTML + fmt.Sprintf(widgetHTMLTemplate,  widgetID, widgetName)
	}
	responseHTML := fmt.Sprintf(homePageHTMLTemplate, widgetsHTML)
	fmt.Fprintf(w, responseHTML)
}

func getWidgetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	widgetStatus, err := q.GetWidgetStatusByID(vars["id"])
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, widgetStatus + "\n")
}

func getWidgetByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	widgetStatus, err := q.GetWidgetStatusByName(vars["name"])
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, widgetStatus + "\n")
}

func setWidgetStatusByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqBody, _ := ioutil.ReadAll(r.Body)
	resp, err := q.SetWidgetStatusByID(vars["id"], string(reqBody))
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, resp + "\n")
}

func setWidgetStatusByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqBody, _ := ioutil.ReadAll(r.Body)
	resp, err := q.SetWidgetStatusByName(vars["name"], string(reqBody))
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, resp + "\n")
}

var q qlcplus.WebsocketConnectionHandler

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Print build information.
	versionString := "qlcplus-http-api"
	versionString = fmt.Sprintf("%s\n  Version: %s %s", versionString, Version, Platform)
	versionString = fmt.Sprintf("%s\n  Commit: %s", versionString, Commit)
	versionString = fmt.Sprintf("%s\n  Built: %s", versionString, BuildDate)
	versionString = fmt.Sprintf("%s\n  Go: %s", versionString, GoVersion)
	log.Println(versionString)

	// Create QLC+ websocket connection handler.
	q = qlcplus.WebsocketConnectionHandler{
		Address: *qlcplusAddress,
	}
	log.Printf("Using QLC+ websocket API on %s\n", *qlcplusAddress)

	// Create HTTP router to server QLC+ API.
	qlcplusRouter := mux.NewRouter().StrictSlash(true)
	qlcplusRouter.HandleFunc("/", homePage)
	qlcplusRouter.HandleFunc("/widgets/id/{id}", setWidgetStatusByID).Methods("POST")
	qlcplusRouter.HandleFunc("/widgets/id/{id}", getWidgetByID)
	qlcplusRouter.HandleFunc("/widgets/name/{name}", setWidgetStatusByName).Methods("POST")
	qlcplusRouter.HandleFunc("/widgets/name/{name}", getWidgetByName)
    log.Printf("Now serving QLC+ HTTP API on %s...\n", *httpAddress)
    log.Fatal(http.ListenAndServe(*httpAddress, qlcplusRouter))
}
