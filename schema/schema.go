package schema

import driver "github.com/sjindal-arista/Arista-QA/memory-driver"

// Driver defines basic func for all storage
type Driver interface {
	AddNewQuestion(newQ string, u *driver.User) error
	SearchQuestion(query string) []*driver.Question
	AnswerQuestion(qID string, ans string, u *driver.User) error
	AddComment(ansID string, comStr string, u *driver.User) error
	AddUser(handle, name, pass string) error
	GetUser(key string) *driver.User
}
