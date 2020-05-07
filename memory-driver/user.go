package driver

import (
	"fmt"
)

type User struct {
	handle   string // handle is the primary key
	fullName string
	passHash string
	asked    []*Question
	answered []*Answer
	stats    *UserStats
}

type UserStats struct {
	reputation int // keep it simple for now. ideally its a function of
	// these things below
	totalAnsViews   int64
	upvotesReceived int
	hoursOnline     int
}

func (d *driver) AddUser(handle, name, pass string) error {
	if _, ok := d.userMapping[handle]; ok {
		return fmt.Errorf("User with handle [%s] exists!", handle)
	}
	newUser := &User{
		handle:   handle,
		fullName: name,
		passHash: pass,
		asked:    []*Question{},
		answered: []*Answer{},
		stats:    &UserStats{},
	}
	d.listUsers = append(d.listUsers, newUser)
	d.userMapping[handle] = newUser
	return nil
}

// returns the user based on primary key - which is user handle
func (d *driver) GetUser(key string) *User {
	res, ok := d.userMapping[key]
	if !ok {
		return nil
	}
	return res
}
