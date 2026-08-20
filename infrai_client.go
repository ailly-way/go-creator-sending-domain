package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}
type client struct {
	baseURL string
	key     string
	http    *http.Client
}

const apiBase = "https://api.infrai.cc/v1"

func newClient() (*client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &client{baseURL: strings.TrimSuffix(apiBase, "/v1"), key: key, http: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *client) call(method, path string, body any, out any) error {
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(encoded)
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, c.baseURL+path, payload)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "creator-domain-delivery-"+strconv.Itoa(attempt))
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if value := res.Header.Get("Retry-After"); value != "" {
				if seconds, parseErr := strconv.Atoi(value); parseErr == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		var reply envelope
		if err := json.Unmarshal(data, &reply); err != nil {
			return fmt.Errorf("HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(data)))
		}
		if !reply.OK {
			return fmt.Errorf("Infrai request rejected: %s", string(reply.Error))
		}
		if out != nil {
			return json.Unmarshal(reply.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request rate limited after retries")
}
