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

type Weddings struct {
}

var (
	MyStats    NewStats
	Error      = ""
	MyWeddings Weddings
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
	// fmt.Print("postBody\n")
	// fmt.Print(string(postBody))
	// fmt.Print("\n")

	callGetStats(postBody)
}

func callGetStats(postBody []byte) {
	var record NewStats
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://128.199.84.236/graphql", responseBody)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Origin", "http://fakewebsite.com")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print("getStat: ")
		fmt.Print(err)
		Error = err.Error()
	}

	if resp != nil {
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
	}

	MyStats = record
}

func GetAllWeddings() {

	var record Weddings
	postBody, _ := json.Marshal(graphql{
		OperationName: "readHaslettWeddingss",
		Query:         "{\n readHaslettWeddingss {\n edges {\n node {\n ID\n Name\n Date\n Package\n }\n }\n }\n}\n",
	})
	responsebody := bytes.NewBuffer(postBody)
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://128.199.84.236/graphql", responsebody)
	if err != nil {
		Error = err.Error()
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("origin", "www.fakewebsite.com")

	resp, err2 := client.Do(req)

	if err2 != nil {
		Error = err2.Error()
	}

	if resp != nil {
		defer resp.Body.Close()
		err3 := json.NewDecoder(resp.Body).Decode(&record)
		if err3 != nil {
			Error = err3.Error()
		}
	}
	fmt.Print(record)
	MyWeddings = record
}
