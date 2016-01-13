// An example that demonstrates the basic operations (Create, Retrieve, Update, Delete)
// Duplicated code exists in order to keep things simple
// May be refactored as long as the example is easy to follow
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Todo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var currentId int
var todos []*Todo

func TodoFind(id int) (*Todo, error) {
	for _, todo := range todos {
		if todo.Id == id {
			return todo, nil
		}
	}

	return nil, fmt.Errorf("Unable to find Todo entry with id %d", id)
}

func TodoDelete(id int) error {
	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Unable to find Todo entry with id %d", id)
}

func TodoNew() *Todo {
	currentId += 1
	var todo Todo = Todo{currentId, "", false}
	todos = append(todos, &todo)
	return &todo
}

func HandlerList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	js, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func HandlerRetrieve(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id int
	var err error
	if id, err = strconv.Atoi(ps.ByName("id")); err != nil {
		panic(err)
	}

	var todo *Todo
	todo, err = TodoFind(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	js, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func HandlerCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var todo *Todo = TodoNew()

	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id int
	var err error
	if id, err = strconv.Atoi(ps.ByName("id")); err != nil {
		panic(err)
	}

	var todo *Todo
	todo, err = TodoFind(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func HandlerDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id int
	var err error
	if id, err = strconv.Atoi(ps.ByName("id")); err != nil {
		panic(err)
	}

	err = TodoDelete(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := httprouter.New()
	router.GET("/todos/", HandlerList)
	router.POST("/todos/", HandlerCreate)
	router.GET("/todos/:id", HandlerRetrieve)
	router.PATCH("/todos/:id", HandlerUpdate)
	router.DELETE("/todos/:id", HandlerDelete)
	log.Fatal(http.ListenAndServe(":9090", router))
}
