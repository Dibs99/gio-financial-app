package apiCalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
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
			MyCategory string      `json:"mycategory"`
			Percentage string      `json:"percentage"`
			Total      string      `json:"total"`
			Colour     color.NRGBA `json:"colour"`
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
	Data struct {
		ReadHaslettWeddingss struct {
			Edges []WeddingNode `json:"edges"`
		} `json:"readHaslettWeddingss"`
	} `json:"data"`
}

type WeddingNode struct {
	Node struct {
		ID                    string `json:"ID"`
		Name                  string `json:"Name"`
		Date                  string `json:"Date"`
		Package               string `json:"Package"`
		Venue                 string `json:"Venue"`
		Aisle                 string `json:"Aisle"`
		Signing               string `json:"Signing"`
		Exit                  string `json:"Exit"`
		Notes                 string `json:"Notes"`
		ReceptionWalkIn       string `json:"ReceptionWalkIn"`
		FirstDance            string `json:"FirstDance"`
		DinnerProvided        bool   `json:"DinnerProvided"`
		HowDidTheyHearAboutMe string `json:"HowDidTheyHearAboutMe"`
		PersonalDetails       string `json:"PersonalDetails"`
		Gross                 string `json:"Gross"`
		Net                   string `json:"Net"`
		Stage                 string `json:"Stage"`
	} `json:"node"`
}

var (
	MyStats    NewStats
	Error      = ""
	MyWeddings Weddings
	requestDestination = "localhost:8888"
)

func GetAllStats() {
	postBody, _ := json.Marshal(graphql{
		OperationName: "HaslettBankStatements",
		Query:         "{readBankStatements{MyCategory\nPercentage\nTotal}}",
	})

	callGetStats(postBody)
}

func GetStatsWithDate(startDate string, endDate string) {
	postBody, _ := json.Marshal(graphql{
		OperationName: "HaslettBankStatements",
		Query:         fmt.Sprintf("{readBankStatements(startDate:\"%v\", endDate:\"%v\"){MyCategory\nPercentage\nTotal}}", startDate, endDate),
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
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v/graphql", requestDestination), responseBody)
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
		// fmt.Print(string(body))
		if err2 != nil {
			log.Fatal("getStat err2: ")
			log.Fatal(err2.Error())
		}
		random := func() uint8 {
			return uint8(rand.Intn(80) + 80)
		}
		// iterate over record and add a colour
		for i := range record.Data.ReadBankStatements {
			colour := color.NRGBA{G: random(), B: random(), R: random(), A: 0xFF}
			record.Data.ReadBankStatements[i].Colour = colour
		}

	}
	MyStats = record
}

func GetAllWeddings() {

	var record Weddings
	postBody, _ := json.Marshal(graphql{
		OperationName: "readHaslettWeddingssConnection",
		Query:         "{ readHaslettWeddingss { edges { node { ID Name Date Package Venue Aisle Signing Exit ReceptionWalkIn FirstDance Notes DinnerProvided HowDidTheyHearAboutMe PersonalDetails Gross Net } } } }",
	})
	responsebody := bytes.NewBuffer(postBody)
	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v/graphql", requestDestination), responsebody)
	if err != nil {
		Error = err.Error()
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("origin", "www.fakewebsite.com")

	resp, err2 := client.Do(req)

	if err2 != nil {
		Error = err2.Error()
	}

	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Print(string(body))
	if resp != nil {
		defer resp.Body.Close()
		err3 := json.NewDecoder(resp.Body).Decode(&record)
		if err3 != nil {
			Error = err3.Error()
		}
	}
	// fmt.Print(record)
	MyWeddings = record
}
