package driver

import (
	"fmt"

	"github.com/s-jindal/QA-Platform/utils"
)

type Answer struct {
	ansID    string
	content  string
	user     *User
	question *Question
	comments []*Comment
	stats    *Stats
}

type Stats struct {
	views     int
	upvotes   int
	downvotes int
}

func NewAnswerStat() *Stats {
	return &Stats{
		views: 1,
	}
}

func NewAnswer(ans string, u *User, q *Question) *Answer {
	return &Answer{
		ansID:    utils.GenerateUUID("ANS"),
		content:  ans,
		user:     u,
		question: q,
		comments: []*Comment{},
		stats:    NewAnswerStat(),
	}
}

func (d *driver) AnswerQuestion(qID string, ans string, u *User) error {
	var quesToBeAns *Question
	for _, ques := range d.listQues {
		if ques.quesID == qID {
			quesToBeAns = ques
			break
		}
	}

	if quesToBeAns != nil {
		newAnswer := NewAnswer(ans, u, quesToBeAns)
		quesToBeAns.answers = append(quesToBeAns.answers, newAnswer)
		return nil
	}
	return fmt.Errorf("Question with QID [%s] not found", qID)
}
