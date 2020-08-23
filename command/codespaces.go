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

	response, err := api.GetCodespaces(apiClient, currentUser)
	if err != nil {
		return err
	}

	for _, codespace := range response.Codespaces {
		codespaceDetails, err := api.GetCodespaceDetails(apiClient, currentUser, codespace.Name)
		if err != nil {
			return err
		}

		hasUnpushedChanges := ""
		if codespaceDetails.Environment.HasUnpushedChanges {
			hasUnpushedChanges = "Has unpushed changes"
		}

		fmt.Printf("%s\t%s\t%s\t%s\n",
			codespace.Name,
			codespaceDetails.Environment.State,
			codespaceDetails.Environment.SkuDisplayName,
			hasUnpushedChanges,
		)
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
