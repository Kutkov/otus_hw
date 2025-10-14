package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DialogClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDialogClient(baseURL string) *DialogClient {
	return &DialogClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type SendMessageRequest struct {
	Text string `json:"text"`
}

type DialogMessageResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func (c *DialogClient) SendMessage(fromUserID, toUserID, text, authToken string) error {
	url := fmt.Sprintf("%s/dialog/%s/send", c.baseURL, toUserID)

	reqBody := SendMessageRequest{Text: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dialog service returned status: %d", resp.StatusCode)
	}

	return nil
}

func (c *DialogClient) GetMessages(userID1, userID2, authToken string) ([]DialogMessageResponse, error) {
	url := fmt.Sprintf("%s/dialog/%s/list", c.baseURL, userID2)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dialog service returned status: %d", resp.StatusCode)
	}

	var messages []DialogMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return messages, nil
}
