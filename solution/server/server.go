package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type data struct {
	Item string `json:"item"`
}

type inputJSON struct {
	ReqType string `json:"req_type"`
	Data    []data `json:"data"`
}

type outputJSONStruct struct {
	ResType string `json:"res_type"`
	Result  string `json:"result"`
	Data    string `json:"data"`
}

type parseError struct {
	err            error
	httpStatusCode int
}

func (pe *parseError) Error() string {
	return fmt.Sprintf("ERROR %v", pe.err)
}

type person struct {
	name     string
	address  string
	location string
}

var statesOrder = [8]string{"AZ", "CA", "ID", "IN", "MA", "OK", "PA", "VA"}

var stateNames = map[string]string{
	"AZ": "Arizona",
	"CA": "California",
	"ID": "Idaho",
	"IN": "Indiana",
	"MA": "Massachusetts",
	"OK": "Oklahoma",
	"PA": "Pennsylvania",
	"VA": "Virginia",
}

func main() {
	port, ok := os.LookupEnv("LISTEN_PORT")
	if !ok {
		fmt.Println("Listened port do not set in env, will used :8080")
	}
	r := mux.NewRouter()
	r.HandleFunc("/", parserHandler)
	http.ListenAndServe(":"+port, r)
}

func parserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		data := outputJSONStruct{}
		data.Result = "fail"
		data.Data = "only POST method allowed"
		out, _ := json.Marshal(data)
		w.Write(out)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data := outputJSONStruct{}
		data.Result = "fail"
		data.Data = "can't read the body"
		out, _ := json.Marshal(data)
		w.Write(out)
		return
	}
	msg := inputJSON{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data := outputJSONStruct{}
		data.Result = "fail"
		data.Data = "bad JSON"
		out, _ := json.Marshal(data)
		w.Write(out)
		return
	}
	out := outputJSONStruct{}
	err = msg.parse(&out)
	if err != nil {
		switch err.(type) {
		case *parseError:
			err := err.(*parseError)
			w.WriteHeader(err.httpStatusCode)
			data := outputJSONStruct{}
			data.Result = "fail"
			data.Data = err.Error()
			out, _ := json.Marshal(data)
			w.Write(out)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			data := outputJSONStruct{}
			data.Result = "fail"
			data.Data = "unknown err type"
			out, _ := json.Marshal(data)
			w.Write(out)
			return
		}

	}

	data, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data := outputJSONStruct{}
		data.Result = "fail"
		data.Data = "can't marshal JSON"
		out, _ := json.Marshal(data)
		w.Write(out)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (ij *inputJSON) parse(out *outputJSONStruct) error {
	if ij.ReqType != "parseAddress" {
		return &parseError{err: errors.New("wrong ReqType"), httpStatusCode: http.StatusBadRequest}
	}

	stateMembers := make(map[string][]person)
	for _, str := range ij.Data {
		info := strings.Split(str.Item, ",")
		if len(info) == 1 {
			if info[0] == "" {
				continue
			}
			return &parseError{err: errors.New("wrong item struct"), httpStatusCode: http.StatusBadRequest}
		}
		if len(info) != 3 {
			return &parseError{err: errors.New("wrong item struct"), httpStatusCode: http.StatusBadRequest}
		}
		for i := range info {
			info[i] = strings.TrimSpace(info[i])
		}
		loc := strings.Split(info[2], " ")
		state := loc[len(loc)-1]

		valideState := false
		for i := range statesOrder {
			if statesOrder[i] == state {
				valideState = true
			}
		}
		if !valideState {
			return &parseError{err: errors.New("wrong state"), httpStatusCode: http.StatusBadRequest}
		}

		newPerson := person{name: info[0],
			address: info[1], location: strings.TrimRight(info[2], (" " + state))}
		members := append(stateMembers[state], newPerson)
		stateMembers[state] = members
	}
	var outData string
	for i, st := range statesOrder {
		data, ok := stateMembers[st]
		if !ok {
			continue
		}
		outData += stateNames[st]
		sort.SliceStable(data, func(i, j int) bool {
			return data[i].name < data[j].name
		})
		for _, member := range data {
			outData += fmt.Sprint("\n..... ", member.name,
				member.address, " ", member.location, " ", stateNames[st])
		}
		if i != len(statesOrder)-1 {
			outData += "\n"
		}
	}
	out.ResType = ij.ReqType
	out.Result = "success"
	out.Data = outData
	return nil
}
