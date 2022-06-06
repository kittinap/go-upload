package service

import (
	"awesomeProject/model"
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
)

type HealthCheckService interface {
	ReadFile(file multipart.File) ([][]string, error)
	Checker(i int, chInput chan string, wg *sync.WaitGroup)
	Request(url string) bool
	Result() model.HealthCheck
}

type healthCheckService struct {
	state model.HealthCheck
}

func NewHealthCheckService() HealthCheckService {
	return &healthCheckService{}
}

func (h *healthCheckService) resetState() {
	h.state.Up = 0
	h.state.Down = 0
}

func (h *healthCheckService) Result() model.HealthCheck {
	return h.state
}

func (h *healthCheckService) ReadFile(file multipart.File) ([][]string, error) {
	h.resetState()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return records, err
	}
	return records, nil
}

func (h *healthCheckService) Checker(i int, chInput chan string, wg *sync.WaitGroup) {
	fmt.Println("worker starting number", i)
	for val := range chInput {
		fmt.Println("worker number :", i, "val :: ", val)
		if h.Request(val) {
			h.state.Up += 1
		} else {
			h.state.Down += 1
		}
	}
	fmt.Println("🔥 end worker :: ", i)
	wg.Done()
}

func (h *healthCheckService) Request(url string) bool {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err) // should be handled properly
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return true
	}
	return false
}
