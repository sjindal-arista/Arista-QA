package elasticsearch

import (
	"flag"

	"github.com/aristanetworks/glog"
)

type matchAll struct{}

// Test -
func Test() {
	flag.Parse()

	eclient := GenerateDriver()

	err := eclient.AddNewQuestion("My first question?", "user1", nil)
	if err != nil {
		glog.Fatal(err)
	}
	err = eclient.AddNewQuestion("My Second question?", "user2", nil)
	if err != nil {
		glog.Fatal(err)
	}

	qid := "ques_1588931754908248389_81"
	err = eclient.AnswerQuestion(qid, "Answer to the 1st qusn", "user3")
	if err != nil {
		glog.Fatal(err)
	}
	err = eclient.AddNewQuestion("My Third question?", "user2", []string{"cvp"})
	if err != nil {
		glog.Fatal(err)
	}

	ansid := "ans_1588932248282659190_81"
	err = eclient.AddComment(ansid, "a good comment", "user1")
	if err != nil {
		glog.Fatal(err)
	}
}
