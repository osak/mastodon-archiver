package mastodon

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"
)

type MastodonClient struct {
	accessToken string
	baseURL     string
}

type Account struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Acct string `json:"acct"`
	DisplayName string `json:"display_name"`
	Locked bool `json:"locked"`
	Bot bool `json:"bot"`
	CreatedAt string `json:"created_at"`
	Note string `json:"note"`
	Uri string `json:"uri"`
}

type Status struct {
	Id string `json:"id"`
	CreatedAt string `json:"created_at"`
	InReplyToId string `json:"in_reply_to_id"`
	Sensitive bool `json:"sensitive"`
	SpoilerText string `json:"spoiler_text"`
	Visibility string `json:"visibility"`
	Language string `json:"language"`
	Uri string `json:"uri"`
	Content string `json:"content"`
	Reblog *Status `json:"reblog"`
	Application string `json:"application"`
	Account Account `json:"account"`
}


func NewMastodonClient(baseURL string, accessToken string) *MastodonClient {
	return &MastodonClient{
		baseURL: baseURL,
		accessToken: accessToken,
	}
}

func (c *MastodonClient) GetHome() error {
	resp, err := c.doGet("/api/v1/timelines/home")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)

	return nil
}

func (c *MastodonClient) GetHomeRaw() ([]string, error) {
	resp, err := c.doGet("/api/v1/timelines/home")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	postBlobs := make([]interface{}, 0)
	err = json.Unmarshal(body, &postBlobs)
	if err != nil {
	    return nil, fmt.Errorf("server returned non-JSON response: %w", err)
	}

	postStrings := make([]string, len(postBlobs))
	for i, post := range postBlobs {
	    postString, err := json.Marshal(post)
	    if err != nil {
		return nil, fmt.Errorf("failed to marshal post: %w", err)
	    }
	    postStrings[i] = string(postString)
	}

	return postStrings, nil
}

func (c *MastodonClient) doGet(path string) (*http.Response, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", c.baseURL + path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer " + c.accessToken)
	req.Header.Add("Content-Type", "application/json")
	
	return httpClient.Do(req)
}

