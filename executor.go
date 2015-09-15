package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "hello!")
  }))

	http.Handle("/", router)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
