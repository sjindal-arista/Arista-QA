package schema

import driver "github.com/s-jindal/QA-Platform/memory-driver"

// Driver defines basic func for all storage
type Driver interface {
	// // return the actual question string
	// GetQuestion(key string) string
	// // return the key of given question
	// GetKey(ques string) string
	// // list of questions which matched
	// // SearchQuestions(query string) []string
	// GetUser(key string) string
	AddNewQuestion(newQ string, u *driver.User) error
}
