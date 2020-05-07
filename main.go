package main

import (
	"fmt"

	driver "github.com/sjindal-arista/Arista-QA/memory-driver"
	"github.com/sjindal-arista/Arista-QA/schema"
)

var handler schema.Driver

// create schemas for question, answer and user objects.
func init() {
	handler = driver.GenerateDriver()
}

func main() {
	fmt.Println("Hello, this is from Q&A team!")
	handler.AddUser("sjindal", "shivam", "")
	handler.AddUser("pranjit_", "pranjit bharali", "")
	handler.AddNewQuestion("Why is TAcc so difficult?", handler.GetUser("sjindal"))
	handler.AddNewQuestion("Why is Cvp more interesting than EOS?", handler.GetUser("pranjit_"))

	matchs := handler.SearchQuestion("tacc difficult")
	for _, m := range matchs {
		fmt.Println(m)
	}
}

// +build ignore

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func main() {
// 	http.HandleFunc("/monkeys", handler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
