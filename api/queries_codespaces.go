package api

import "fmt"

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

// CodespaceDetails represents details about a codespace.
type CodespaceDetails struct {
	Name        string
	GUID        string
	State       string
	URL         string
	TokenURL    string `json:"token_url"`
	Environment struct {
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
