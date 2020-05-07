package driver

import (
	"fmt"

	"github.com/sjindal-arista/Arista-QA/utils"
)

type Comment struct {
	cID     string
	comment string
	user    *User
}

func NewComment(content string, u *User) *Comment {
	return &Comment{
		cID:     utils.GenerateUUID("COMMENT"),
		comment: content,
		user:    u,
	}
}

func (d *driver) AddComment(ansID string, comStr string, u *User) error {
	var a *Answer
	for _, ans := range d.listAnswers {
		if ans.ansID == ansID {
			a = ans
			break
		}
	}
	if a != nil {
		a.comments = append(a.comments, NewComment(comStr, u))
		return nil
	}
	return fmt.Errorf("Answer with ID [%s] not found", ansID)
}
