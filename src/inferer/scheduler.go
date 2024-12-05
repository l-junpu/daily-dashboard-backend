package inferer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Endpoint handling API done in tcp-api.go
type Endpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Scheduler struct {
	MaxRetries    int        `json:"maxRetries"`
	RetryInterval int        `json:"retryInterval"`
	Endpoints     []Endpoint `json:"endpoints"`
}

func CreateScheduler(cwd string, err error) (*Scheduler, error) {
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve cwd in main: %w", err)
	}

	scheduler := Scheduler{
		Endpoints: make([]Endpoint, 0),
	}
	scheduler.IntitializeEndpoints(filepath.Join(cwd, "cfg", "endpoints.json"))

	return &scheduler, nil
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

	if err := json.Unmarshal(byteArr, b); err != nil {
		log.Fatal("error unmarshalling byte array into endpoints array: %w", err)
	}
}

func (b *Scheduler) HasAvailableEndpoint() bool {
	return len(b.Endpoints) > 0
}

func (b *Scheduler) GetEndpoint() (*Endpoint, error) {
	maxRetries := b.MaxRetries
	for maxRetries >= 0 {
		if !b.HasAvailableEndpoint() {
			time.Sleep(time.Duration(b.RetryInterval) * time.Second)
			maxRetries -= 1
			continue
		}
		Endpoint := b.Endpoints[len(b.Endpoints)-1]
		b.Endpoints = b.Endpoints[:len(b.Endpoints)-1]
		return &Endpoint, nil
	}
	return nil, fmt.Errorf("there are no available schedulers at the moment")
}

func (b *Scheduler) ReturnEndpoint(endPt *Endpoint) {
	b.Endpoints = append(b.Endpoints, *endPt)
}
