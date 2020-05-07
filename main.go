package main

import (
	"fmt"

	driver "github.com/s-jindal/QA-Platform/memory-driver"
	"github.com/s-jindal/QA-Platform/schema"
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
