package hooks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Hook struct {
	Name       string
	Conditions []Condition
	Actions    []Action
}

type Condition struct {
	SupervisorHealthy bool
	// Add more condition fields as needed
}

type Action struct {
	Name               string
	Kind               string
	URL                string
	Method             string
	ExpectedStatusCode int
	// Add more action fields as needed
}

func NewHook(data map[string]interface{}) (*Hook, error) {
	var hook Hook
	err := mapstructure.Decode(data, &hook)
	if err != nil {
		return nil, err
	}
	return &hook, nil
}

func ExecuteHTTPAction(action Action) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(action.Method, action.URL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != action.ExpectedStatusCode {
		return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, action.ExpectedStatusCode)
	}

	return nil
}
