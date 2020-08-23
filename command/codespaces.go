package command

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(codespacesCmd)

	codespacesCmd.AddCommand(codespacesListCmd)
}

func codespacesList(cmd *cobra.Command, args []string) error {
	ctx := contextForCommand(cmd)
	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		return err
	}

	currentUser, err := api.CurrentLoginName(apiClient)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("vscs_internal/user/%s/codespaces", currentUser)

	type Codespace struct {
		Name     string `json:"name"`
		GUID     string `json:"guid"`
		State    string `json:"state"`
		URL      string `json:"url"`
		TokenURL string `json:"token_url"`
	}

	type Response struct {
		Codespaces []Codespace `json:"codespaces"`
	}

	var response Response

	err = apiClient.REST("GET", endpoint, nil, &response)
	if err != nil {
		return err
	}

	for _, v := range response.Codespaces {
		fmt.Println(v.Name)
	}

	return nil
}

var codespacesCmd = &cobra.Command{
	Use:   "codespaces <command>",
	Short: "Create, list, and delete codespaces",
	Long:  `Work with GitHub codespaces`,
	Example: heredoc.Doc(`
	$ gh codespaces create vsls-contrib/guestbook
	$ gh codespaces delete srivatsn-vsls-contrib-asd23
	$ gh codespaces list
	`),
}

var codespacesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List codespaces for current user",
	Args:  cmdutil.NoArgsQuoteReminder,
	Example: heredoc.Doc(`
	$ gh codespaces list
	`),
	RunE: codespacesList,
}
