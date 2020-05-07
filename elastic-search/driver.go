package elasticsearch

import (
	"context"
	"fmt"

	"github.com/aristanetworks/glog"
	"github.com/sjindal-arista/Arista-QA/schema"
	"github.com/sjindal-arista/Arista-QA/utils"
)

// ElasticDriver -
type ElasticDriver struct {
	*Client
}

// GenerateDriver -
func GenerateDriver() *ElasticDriver {
	ec := NewClient("http://localhost:9200")

	err := initialise(ec)
	if err != nil {
		glog.Fatal(err)
		return nil
	}
	return &ElasticDriver{
		Client: ec,
	}
}

func initialise(eclient *Client) error {
	ctx := context.Background()
	glog.Infof("Creating index...")
	err := eclient.CreateIndex(ctx, "user", nil, 400)
	if err != nil {
		glog.Errorf("Error in creating question index: %v", err)
		return err
	}
	err = eclient.CreateIndex(ctx, "question", nil, 400)
	if err != nil {
		glog.Errorf("Error in creating question index: %v", err)
		return err
	}
	err = eclient.CreateIndex(ctx, "answer", nil, 400)
	if err != nil {
		glog.Errorf("Error in creating answer index: %v", err)
		return err
	}
	return nil
}

type questionDoc struct {
	key  string
	qusn string
	user string
	tags []string
}

type answerDoc struct {
	key      string
	qusnID   string
	ans      string
	user     string
	stats    schema.Stats
	comments []schema.Comment
}

// AddNewQuestion -
func (driver *ElasticDriver) AddNewQuestion(newQ, user string, tags []string) error {
	id := utils.GenerateUUID("ques")
	docs := []BulkRequest{
		{
			Index: "question",
			ID:    id,
			Body: map[string]interface{}{
				"data": questionDoc{
					key:  id,
					qusn: newQ,
					user: user,
					tags: tags,
				},
			},
		},
	}

	glog.Infof("Adding question ...")
	err := driver.Bulk(context.Background(), docs)
	if err != nil {
		glog.Errorf("Error in adding question: %v", err)
		return fmt.Errorf("Error in adding question: %v", err)
	}
	glog.Infof("Question added")
	return nil
}

// SearchQuestion -
func (driver *ElasticDriver) SearchQuestion(query string, tags []string) error {
	return fmt.Errorf("Not implememted")
}

// AnswerQuestion -
func (driver *ElasticDriver) AnswerQuestion(qID, ans, user string) error {
	id := utils.GenerateUUID("ans")
	docs := []BulkRequest{
		{
			Index: "answer",
			ID:    id,
			Body: map[string]interface{}{
				"data": answerDoc{
					key:      id,
					qusnID:   qID,
					user:     user,
					stats:    schema.Stats{},
					comments: []schema.Comment{},
				},
			},
		},
	}

	glog.Infof("Adding answer ...")
	err := driver.Bulk(context.Background(), docs)
	if err != nil {
		glog.Errorf("Error in adding answer: %v", err)
		return fmt.Errorf("Error in adding answer: %v", err)
	}
	glog.Infof("Answer added")
	return nil
}

//AddComment -
func (driver *ElasticDriver) AddComment(ansID, comStr, user string) error {
	return fmt.Errorf("Not implememted")
}
