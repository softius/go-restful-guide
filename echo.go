// Just a simple echo operation that returns the same string that is given to it.
// curl 'http://localhost:9090/echo/Hello' -v -H 'Accept: application/json'
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Message struct {
	Message string `json:"message"`
}

func echo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message := Message{ps.ByName("message")}

	js, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	router := httprouter.New()
	router.GET("/echo/:message", echo)
	log.Fatal(http.ListenAndServe(":9090", router))
}
