package driver

import "fmt"

type Question struct {
	q       string
	user    *User
	answers []*Answer
}

func (d *driver) AddNewQuestion(newQ string, u *User) error {
	// assert user already exists
	if !d.userExists(u) {
		return fmt.Errorf("User does not exist")
	}
	questionObj := &Question{
		q:    newQ,
		user: u,
	}
	d.listQues = append(d.listQues, questionObj)
	return nil
}
