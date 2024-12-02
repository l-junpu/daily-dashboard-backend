package inferer

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Endpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Scheduler struct {
	Endpoints []Endpoint
}

func (b *Scheduler) IntitializeEndpoints(EndpointConfigPath string) {
	jsonFile, err := os.Open(EndpointConfigPath)
	if err != nil {
		log.Fatal("unable to open endpoint config file: %w", err)
	}
	defer jsonFile.Close()

	byteArr, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("unable to parse json file into byte array: %w", err)
	}

	if err := json.Unmarshal(byteArr, &b.Endpoints); err != nil {
		log.Fatal("error unmarshalling byte array into endpoints array: %w", err)
	}
}
