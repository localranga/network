package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/localranga/tttt/request"
)

type NetworkPipeline struct {
	outputs map[string]chan<- *request.ResponseData
}

func NewNetworkPipeline() *NetworkPipeline {
	return &NetworkPipeline{
		outputs: make(map[string]chan<- *request.ResponseData),
	}
}

func (n *NetworkPipeline) Output(apiID string) <-chan *request.ResponseData {
	output := make(chan *request.ResponseData)
	n.outputs[apiID] = output
	return output
}

func (n *NetworkPipeline) Start() {
	ticker := time.NewTicker(1 * time.Second) // Adjust the interval according to your needs
	defer ticker.Stop()

	for range ticker.C {
		for apiID, output := range n.outputs {
			response, err := fetchData(apiID)
			if err != nil {
				fmt.Printf("Error fetching data for %s: %v\n", apiID, err)
				continue
			}

			output <- response
		}
	}
}

func fetchData(apiID string) (*request.ResponseData, error) {
	apiData, err := ioutil.ReadFile("requests/" + apiID + ".go")
	if err != nil {
		return nil, err
	}

	data, err := request.FetchDataFromAPI(apiData)
	if err != nil {
		return nil, err
	}

	response := &request.ResponseData{
		Timestamp: time.Now(),
		Data:      data,
	}

	return response, nil
}