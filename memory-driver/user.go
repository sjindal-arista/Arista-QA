package driver

import (
	"fmt"
)

type User struct {
	Handle   string // Handle is the primary key
	FullName string
	PassHash string
	Asked    []*Question
	Answered []*Answer
	Stats    *UserStats
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
		return fmt.Errorf("User with Handle [%s] exists!", handle)
	}
	newUser := &User{
		Handle:   handle,
		FullName: name,
		PassHash: pass,
		Asked:    []*Question{},
		Answered: []*Answer{},
		Stats:    &UserStats{},
	}
	d.listUsers = append(d.listUsers, newUser)
	d.userMapping[handle] = newUser
	return nil
}

// returns the user based on primary key - which is user Handle
func (d *driver) GetUser(key string) *User {
	res, ok := d.userMapping[key]
	if !ok {
		return nil
	}
	return res
}
