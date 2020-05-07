package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	driver "github.com/sjindal-arista/Arista-QA/memory-driver"
	"github.com/sjindal-arista/Arista-QA/schema"
)

var handler schema.Driver

// create schemas for question, answer and user objects.
func init() {
	handler = driver.GenerateDriver()
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is from Q&A team!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homepage)
	router.HandleFunc("/QA/adduser", addUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func addUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, " 'Error': Error in decoding body: %v", err)
	} else {
		var user driver.User
		json.Unmarshal(reqBody, &user)
		if err := handler.AddUser(user.Handle, user.FullName, user.PassHash); err != nil {
			fmt.Fprintf(w, " 'Error': could not add user %v", err)
		}
	}
}
