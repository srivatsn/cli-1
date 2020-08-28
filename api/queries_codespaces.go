package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cli/cli/internal/ghrepo"
)

// Codespace represents a single codespace.
type Codespace struct {
	Name     string
	GUID     string
	State    string
	URL      string
	TokenURL string `json:"token_url"`
}

// Codespaces represents the response from the codespaces list api.
type Codespaces struct {
	Codespaces []Codespace
}

// Environment represents the in the VSCS service.
type Environment struct {
	ID           string
	Type         string
	FriendlyName string
	State        string
	Seed         struct {
		Type      string
		Moniker   string
		GitConfig struct {
			UserName  string
			UserEmail string
		}
	}
	Connection struct {
		SessionID   string
		SessionPath string
		ServiceURI  string
	}
	RecentFolders            []string
	Location                 string
	PlanID                   string
	AutoShutdownDelayMinutes int
	SkuName                  string
	SkuDisplayName           string
	LastStateUpdateReason    string
	HasUnpushedChanges       bool
}

// CodespaceDetails represents details about a codespace.
type CodespaceDetails struct {
	Name        string
	GUID        string
	State       string
	URL         string
	TokenURL    string `json:"token_url"`
	Environment Environment
}

// CodespaceCreateRequest represents the request body for a codespace create request
type CodespaceCreateRequest struct {
	RepositoryID int    `json:"repository_id"`
	Ref          string `json:"ref,omitempty"`
	PullRequest  struct {
		RepositoryID      int `json:"repository_id"`
		PullRequestNumber int `json:"pull_request_number"`
	} `json:"pull_request,omitempty"`
	Location string `json:"location,omitempty"`
	Sku      string `json:"sku,omitempty"`
}

// CodespaceCreateResponse is the response from a codespace create request
type CodespaceCreateResponse struct {
	GUID     string
	Name     string
	State    string
	URL      string
	TokenURL string `json:"token_url"`
}

// GetCodespaces gets the codespaces for the given user.
func GetCodespaces(client *Client, currentUsername string) (*Codespaces, error) {
	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces", currentUsername)

	var response Codespaces
	err := client.REST("GET", endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCodespaceDetails gets the details of the given codespace name.
func GetCodespaceDetails(client *Client, currentUsername string, codespaceName string) (*CodespaceDetails, error) {
	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces/%s", currentUsername, codespaceName)

	var details CodespaceDetails
	err := client.REST("GET", endpoint, nil, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

// GetCodespaceToken returns the token for that codespaces to talk to the VSCS API.
func GetCodespaceToken(client *Client, currentUsername string, codespaceName string) (*string, error) {
	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces/%s/token", currentUsername, codespaceName)

	var response struct {
		Token string
	}
	err := client.REST("POST", endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Token, nil
}

// StartCodespace resumes a suspended codespace.
func StartCodespace(client *Client, currentUsername string, codespaceName string) error {
	codespace, err := GetCodespaceDetails(client, currentUsername, codespaceName)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("vscs_internal/proxy/environments/%s/start", codespace.Environment.ID)
	var response Environment
	err = client.REST("POST", endpoint, nil, &response)
	if err != nil {
		return err
	}

	return nil
}

// SuspendCodespace suspends a codespace.
func SuspendCodespace(client *Client, currentUsername string, codespaceName string) error {
	codespace, err := GetCodespaceDetails(client, currentUsername, codespaceName)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://online.visualstudio.com/api/v1/Environments/%s/shutdown", codespace.Environment.ID)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	token, err := GetCodespaceToken(client, currentUsername, codespaceName)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	response, err := client.http.Do(req)
	if err != nil {
		return err
	}

	success := response.StatusCode >= 200 && response.StatusCode < 300
	if !success {
		return handleHTTPError(response)
	}

	return nil
}

// DeleteCodespace deletes a codespace.
func DeleteCodespace(client *Client, currentUsername string, codespaceName string) error {
	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces/%s", currentUsername, codespaceName)

	err := client.REST("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateCodespace creates a codespace.
func CreateCodespace(client *Client, currentUsername string, repoName string) (*string, error) {
	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces", currentUsername)

	repoID, err := getRepoID(client, repoName)
	if err != nil {
		return nil, err
	}

	location, err := getCurrentLocation(client)
	if err != nil {
		return nil, err
	}

	body := &CodespaceCreateRequest{
		RepositoryID: repoID,
		Location:     location,
		Sku:          "basicLinux",
	}

	requestByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	requestBody := bytes.NewReader(requestByte)

	var response CodespaceCreateResponse
	err = client.REST("POST", endpoint, requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Name, nil
}

func getRepoID(client *Client, repoName string) (int, error) {
	repo, err := ghrepo.FromFullName(repoName)
	if err != nil {
		return -1, err
	}

	endpoint := fmt.Sprintf("repos/%s/%s", repo.RepoOwner(), repo.RepoName())
	var response struct {
		ID int
	}
	err = client.REST("GET", endpoint, nil, &response)
	if err != nil {
		return -1, err
	}

	return response.ID, err
}

func getCurrentLocation(client *Client) (string, error) {
	endpoint := "https://online.visualstudio.com/api/v1/locations"
	var response struct {
		Current string
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &response)

	return response.Current, nil
}
