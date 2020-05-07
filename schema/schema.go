package schema

// Question -
type Question struct {
	q       string
	quesID  string
	user    string
	tags    []string
	answers []*Answer
}

// Answer -
type Answer struct {
	ansID    string
	content  string
	user     string
	quesID   string
	comments []*Comment
	stats    *Stats
}

// Stats -
type Stats struct {
	views     int
	upvotes   int
	downvotes int
}

// Comment -
type Comment struct {
	comment string
	user    string
}

// Driver defines basic func for all storage
type Driver interface {
	AddNewQuestion(newQ, user string, tags []string) error
	SearchQuestion(query string, tags []string) []*Question
	AnswerQuestion(qID, ans, user string) error
	AddComment(ansID, comStr, user string) error
	//	AddUser(handle, name, pass string) error
	//	GetUser(key string) *driver.User
}
