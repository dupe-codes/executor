package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	cr "github.com/njdup/executor/coderunner"
)

type Response struct {
	Status int         `bson:"status" json:"status"`
	Data   interface{} `bson:"data" json:"data"`
	Error  error       `bson:"error" json:"error"`
}

func configureRoutes(router *mux.Router) {
	router.Handle("/run", http.HandlerFunc(handleCodeRun)).Methods("POST")
	router.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "hello!")
	}))
	//router.Handle("/health", getHealthStatus).Methods("GET")
}

func sendResponse(rw http.ResponseWriter, response Response) {
	resContent, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(rw, "Error preparing response.", http.StatusInternalServerError)
		return
	}

	if response.Error != nil {
		http.Error(rw, string(resContent), response.Status)
	} else {
		fmt.Fprintf(rw, string(resContent))
	}
}

func writeOutput(cmdOutput string) error {
	// For now, print result. Will need to write to kafka later
	fmt.Fprintln(os.Stdout, "Command run and resulted in the following output:\n", cmdOutput)
	return nil
}

// Expects the incoming request to have the following POST params:
// - name (string name of the code command to be run)
// - code (string of the code to be run)
// - language (string of languag in which code is written)
// - phonenumber (string number of user running code) (TODO)
func handleCodeRun(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	// TODO: Check if needed params present
	runner := cr.CodeRun{
		req.FormValue("name"),
		req.FormValue("code"),
		req.FormValue("language"),
	}

	output, err := runner.Run()
	if err != nil {
    fmt.Println("Error running code...")
		res := Response{http.StatusInternalServerError, nil, err}
		sendResponse(rw, res)
		return
	}

	// Otherwise, code was run successfully
	// Need to somehow write to kafka - maybe POST to a producer server
	// running back on main host?
	err = writeOutput(output)
	var res Response
	if err != nil {
		res = Response{http.StatusInternalServerError, nil, err}
	} else {
		res = Response{http.StatusOK, "Command successfully run", nil}
	}
	sendResponse(rw, res)
}

func main() {
	port := os.Getenv("CODERUNNER_PORT")
	if port == "" {
		port = ":8001" // Default to port 8001
	}

	router := mux.NewRouter()
	configureRoutes(router)
	http.Handle("/", router)

	fmt.Println("Code running server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux)))
}
