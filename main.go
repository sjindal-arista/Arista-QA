package main

import "fmt"
import driver "github.com/s-jindal/QA-Platform/memory-driver"
import "github.com/s-jindal/QA-Platform/schema"

var handler schema.Driver

// create schemas for question, answer and user objects.
func init() {
	handler = driver.GenerateDriver()
}

func main() {
	fmt.Println("Hello, this is Q&A")
}
