package apiCalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var jsonString string = `{
    "data": [
        {
            "category" : "food",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "accommodation",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "travel",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        },
        {
            "category" : "gifts",
            "percentage" : 25,
            "total" : 300,
            "colour" : ""
        }
    ]
}`

type NewStats struct {
	Data struct {
		ReadBankStatements []struct {
			Category   string      `json:"Category"`
			Percentage string      `json:"percentage"`
			Total      float32     `json:"total"`
			Colour     interface{} `json:"colour"`
		} `json:"readBankStatements"`
	} `json:"data"`
}

type graphql struct {
	OperationName string   `json:"operationName"`
	Query         string   `json:"query"`
	Variables     variable `json:"variables"`
}
type variable struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

var (
	MyStats NewStats
)

func GetAllStats() {

	postBody, _ := json.Marshal(graphql{
		OperationName: "HaslettBankStatements",
		Query:         "{readBankStatements{Category\nPercentage\nTotal}}",
	})

	callGetStats(postBody)

}

func GetStatsWithDate(startDate string, endDate string) {
	postBody, _ := json.Marshal(graphql{
		OperationName: "HaslettBankStatements",
		Query:         fmt.Sprintf("{readBankStatements(startDate:\"%v\", endDate:\"%v\"){Category\nPercentage\nTotal}}", startDate, endDate),
		Variables: variable{
			StartDate: startDate,
			EndDate:   endDate,
		},
	})
	fmt.Print("postBody\n")
	fmt.Print(string(postBody))
	fmt.Print("\n")

	callGetStats(postBody)
}

func callGetStats(postBody []byte) {
	var record NewStats
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8888/graphql", responseBody)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Origin", "http://fakewebsite.com")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("getStat: ")
		log.Fatal(err)
	}

	defer resp.Body.Close()

	err2 := json.NewDecoder(resp.Body).Decode(&record)
	// body, err3 := ioutil.ReadAll(resp.Body)
	// if err3 != nil {
	// 	log.Fatal("getStat err3: ")
	// 	log.Fatal(err3)
	// }
	if err2 != nil {
		log.Fatal("getStat err2: ")
		log.Fatal(err2)
	}

	MyStats = record
}
