package driver

type driver struct {
	listQues    []*Question
	listAnswers []*Answer
	listUsers   []*User
}

func GenerateDriver() *driver {
	return &driver{
		listQues:    []*Question{},
		listAnswers: []*Answer{},
		listUsers:   []*User{},
	}
}

func (d *driver) userExists(u *User) bool {
	for _, user := range d.listUsers {
		if user.handle == u.handle {
			return true
		}
	}
	return false
}
