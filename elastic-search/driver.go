package elasticsearch

import (
	"context"
	"encoding/json"
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
	glog.Infof("Creating indices...")
	data := map[string]interface{}{}
	err := eclient.CreateIndex(ctx, "question", data, 400)
	if err != nil {
		glog.Errorf("Error in creating question index: %v", err)
		return err
	}
	err = eclient.CreateIndex(ctx, "answer", data, 400)
	if err != nil {
		glog.Errorf("Error in creating answer index: %v", err)
		return err
	}
	return nil
}

type questionDoc struct {
	Key  string   `json:"key"`
	Qusn string   `json:"qusn"`
	User string   `json:"user"`
	Tags []string `json:"tags"`
}

type answerDoc struct {
	Key      string           `json:"key"`
	QusnID   string           `json:"qusn_id"`
	Ans      string           `json:"ans"`
	User     string           `json:"user"`
	Stats    schema.Stats     `json:"stats"`
	Comments []schema.Comment `json:"comments"`
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
					Key:  id,
					Qusn: newQ,
					User: user,
					Tags: tags,
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
					Key:      id,
					QusnID:   qID,
					Ans:      ans,
					User:     user,
					Stats:    schema.Stats{},
					Comments: []schema.Comment{},
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
	data := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.data.comments.add(params.comment)",
			"params": map[string]interface{}{
				"comment": map[string]string{
					"comment": comStr,
					"user":    user,
				},
			},
		},
	}
	glog.Infof("Adding comment ...")
	res, err := driver.UpdateDoc(context.Background(), "answer", ansID, data)
	if err != nil {
		glog.Errorf("Error in adding comment: %v", err)
		return fmt.Errorf("Error in adding comment: %v", err)
	}
	var resp map[string]interface{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		glog.Errorf("Error in unmarshaling updateDoc response: %v", err)
	}
	glog.Infof("Response is: %v", resp)
	if err, ok := resp["error"]; ok {
		return fmt.Errorf("Error in adding comment: %v", err)
	}
	glog.Infof("Comment added")
	return nil
}
