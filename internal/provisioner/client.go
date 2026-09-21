package provisioner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrUnavailable = errors.New("provisioner unavailable")

type CreateRequest struct {
	Name   string `json:"name"`
	Engine string `json:"engine"`
	SizeGB int32  `json:"sizeGB"`
}

type Database struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Engine   string `json:"engine"`
	SizeGB   int32  `json:"sizeGB"`
	State    string `json:"state"`
	Endpoint string `json:"endpoint,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Create(ctx context.Context, input CreateRequest) (Database, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return Database{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/databases", bytes.NewReader(body))
	if err != nil {
		return Database{}, err
	}

	request.Header.Set("Content-type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Database{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusServiceUnavailable {
		return Database{}, ErrUnavailable
	}

	if response.StatusCode != http.StatusCreated {
		return Database{}, fmt.Errorf("provisioner returned status %d", response.StatusCode)
	}

	var database Database

	if err := json.NewDecoder(response.Body).Decode(&database); err != nil {
		return Database{}, err
	}

	return database, nil
}

func (c *Client) Get(ctx context.Context, id string) (Database, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/databases/"+id, nil)
	if err != nil {
		return Database{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Database{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Database{}, fmt.Errorf("provisioner returned status %d", response.StatusCode)
	}

	var database Database

	if err := json.NewDecoder(response.Body).Decode(&database); err != nil {
		return Database{}, err
	}

	return database, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/databases/"+id, nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent &&
		response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("provisioner returned status %d", response.StatusCode)
	}

	return nil
}
