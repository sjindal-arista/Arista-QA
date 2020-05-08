package elasticsearch

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/aristanetworks/glog"
)

type matchAll struct{}

// Test -
func Test() {
	flag.Parse()

	eclient := NewClient("http://localhost:9200")

	ctx := context.Background()

	body := map[string]interface{}{
		"first_field":  "testindexbody1",
		"second_field": "testindexbody2",
	}

	glog.Infof("Creating index...")
	err := eclient.CreateIndex(ctx, "personal_info", body, 400)

	if err != nil {
		glog.Errorf("Error in creating index: %v", err)
	}

	testData := []BulkRequest{
		{
			Delete: false,
			Index:  "personal_info",
			ID:     "user1",
			Body:   map[string]interface{}{"name": "pranjit", "mobile": 12345},
		},
		{
			Delete: false,
			Index:  "personal_info",
			ID:     "user2",
			Body:   map[string]interface{}{"name": "shivam", "mobile": 67890},
		},
	}

	glog.Infof("Pushing Bulk...")
	err = eclient.Bulk(ctx, testData)
	if err != nil {
		glog.Errorf("Error in sending bulk data: %v", err)
	}

	glog.Infof("Searching data.....")
	body = map[string]interface{}{
		"query": map[string]interface{}{"match_all": matchAll{}},
	}
	res, err := eclient.Search(ctx, []string{"personal_info"}, body, 400)
	//resp := make(map[string]interface{})
	var resp interface{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		glog.Errorf("Error in searching data: %v", err)
	}
	glog.Infof("Response is: %v", resp)
}
