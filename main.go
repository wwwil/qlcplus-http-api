package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/wwwil/qlcplus-websocket-http/pkg/qlcplus"

)

var qlcplusAddress = flag.String("qlcplus", "localhost:9999", "Address of QLC+ websocket")
var httpAddress = flag.String("port", "localhost:8888", "Address for QLC+ HTTP API")

func homePage(w http.ResponseWriter, r *http.Request){
    widgetsMap, err := q.GetWidgetsMap()
    if err != nil {
        log.Println(err)
    }
    response := ""
    for widgetID, widgetName := range widgetsMap {
        response = fmt.Sprintf("%s\n%s - %s", response, widgetID, widgetName)
    }
    // TODO: Print this as a table.
    //response := "<html><tr><th>Widget ID</th><th>Widget Name</th></td>"
    //for widgetID, widgetName := range widgetsMap {
    //    response = fmt.Sprintf("%s<tr><th>%s</th><th>%s</th></tr>", response, widgetID, widgetName)
    //}
    //response = fmt.Sprintf("%s</html>", response)
    fmt.Fprintf(w, response)
}

func getWidgetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    widgetStatus, err := q.GetWidgetStatusByID(vars["id"])
    if err != nil {
        log.Println(err)
    }
    fmt.Fprintf(w, widgetStatus)
}

func getWidgetByName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    widgetStatus, err := q.GetWidgetStatusByName(vars["name"])
    if err != nil {
        log.Println(err)
    }
    fmt.Fprintf(w, widgetStatus)
}

func setWidgetStatusByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reqBody, _ := ioutil.ReadAll(r.Body)
    resp, err := q.SetWidgetStatusByID(vars["id"], string(reqBody))
    if err != nil {
        log.Println(err)
    }
    fmt.Fprintf(w, resp)
}

func setWidgetStatusByName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reqBody, _ := ioutil.ReadAll(r.Body)
    resp, err := q.SetWidgetStatusByName(vars["name"], string(reqBody))
    if err != nil {
        log.Println(err)
    }
    fmt.Fprintf(w, resp)
}

var q qlcplus.Connection

func main() {
    flag.Parse()
    log.SetFlags(0)

    q = qlcplus.Connection{
        Address: *qlcplusAddress,
    }

    qlcplusRouter := mux.NewRouter().StrictSlash(true)
    qlcplusRouter.HandleFunc("/", homePage)
    qlcplusRouter.HandleFunc("/id/{id}", setWidgetStatusByID).Methods("POST")
    qlcplusRouter.HandleFunc("/id/{id}", getWidgetByID)
    qlcplusRouter.HandleFunc("/name/{name}", setWidgetStatusByName).Methods("POST")
    qlcplusRouter.HandleFunc("/name/{name}", getWidgetByName)
    log.Fatal(http.ListenAndServe(*httpAddress, qlcplusRouter))
}