package main

//used just for simple tests
import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	reqBody := bytes.NewBuffer([]byte(testReq))
	req, err := http.NewRequest("POST", "http://localhost:8000/", reqBody)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

var testReq string = `{
	"req_type": "parseAddress",
	"data": [
		{
			"item":"John Daggett, 341 King Road, Plymouth MA"
		},
		{
			"item":"John Daggett, 341 King Road, Plymouth MA"
		},
		{
			"item":"Anthony Daggett, 341 King Road, Plymouth MA"
		},
	  	{
			"item":"Alice Ford, 22 East Broadway, Richmond VA"
	  	},
	  	{
			"item": "Terry Kalkas, 402 Lans Road, Beaver Falls PA"
	  	},
	  	{
			"item": " Eric Adams, 20 Post Road, Sudbury MA"
	  	},
	  	{
			"item": "Sal Carpenter, 73 6th Street, Boston MA"
	  	}
	]
}`
