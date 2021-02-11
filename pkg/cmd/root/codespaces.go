package command

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(codespacesCmd)

	codespacesCmd.AddCommand(codespacesListCmd)
	codespacesCmd.AddCommand(codespacesSuspendCmd)
	codespacesCmd.AddCommand(codespacesResumeCmd)
	codespacesCmd.AddCommand(codespacesDeleteCmd)
	codespacesCmd.AddCommand(codespacesCreateCmd)
	codespacesCreateCmd.Flags().StringP("ref", "r", "", "A ref in the repo from which the codespace will be created. The default branch of the repo is used otherwise")
	codespacesCreateCmd.Flags().StringP("sku", "s", "", "The sku of the codespace (eg: basicLinux)")

}

func getAPIClientAndCurrentUser(cmd *cobra.Command) (*api.Client, string, error) {
	ctx := contextForCommand(cmd)
	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		return nil, "", err
	}

	currentUser, err := api.CurrentLoginName(apiClient)
	if err != nil {
		return nil, "", err
	}

	return apiClient, currentUser, nil
}

func codespacesList(cmd *cobra.Command, args []string) error {
	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
	if err != nil {
		return err
	}

	response, err := api.GetCodespaces(apiClient, currentUser)
	if err != nil {
		return err
	}

	table := utils.NewTablePrinter(cmd.OutOrStdout())

	for _, codespace := range response.Codespaces {
		codespaceDetails, err := api.GetCodespaceDetails(apiClient, currentUser, codespace.Name)
		if err != nil {
			return err
		}

		hasUnpushedChanges := ""
		if codespaceDetails.Environment.HasUnpushedChanges {
			hasUnpushedChanges = "Has unpushed changes"
		}

		table.AddField(codespace.Name, nil, colorfuncForState(codespaceDetails.Environment.State))
		table.AddField(codespaceDetails.Environment.State, nil, colorfuncForState(codespaceDetails.Environment.State))
		table.AddField(codespaceDetails.Environment.SkuDisplayName, nil, nil)
		table.AddField(codespaceDetails.Environment.Seed.Moniker, nil, utils.Blue)
		table.AddField(hasUnpushedChanges, nil, utils.Red)
		table.EndRow()
	}

	err = table.Render()
	if err != nil {
		return err
	}

	return nil
}

func codespacesSuspend(cmd *cobra.Command, args []string) error {
	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
	if err != nil {
		return err
	}

	codespaceName := args[0]
	err = api.SuspendCodespace(apiClient, currentUser, codespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully suspended.\n"), codespaceName)
	return nil
}

func codespacesResume(cmd *cobra.Command, args []string) error {
	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
	if err != nil {
		return err
	}

	codespaceName := args[0]
	err = api.StartCodespace(apiClient, currentUser, codespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully resumed.\n"), codespaceName)
	return nil
}

func codespacesDelete(cmd *cobra.Command, args []string) error {
	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
	if err != nil {
		return err
	}

	codespaceName := args[0]
	err = api.DeleteCodespace(apiClient, currentUser, codespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully deleted.\n"), codespaceName)
	return nil
}

func codespacesCreate(cmd *cobra.Command, args []string) error {
	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
	if err != nil {
		return err
	}

	repoName := args[0]

	ref, err := cmd.Flags().GetString("ref")
	if err != nil {
		return err
	}

	sku, err := cmd.Flags().GetString("sku")
	if err != nil {
		return err
	}

	codespaceName, err := api.CreateCodespace(apiClient, currentUser, repoName, ref, sku)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully created.\n"), codespaceName)
	return nil
}

func colorfuncForState(state string) func(string) string {
	switch state {
	case "Available":
		return utils.Green
	case "Shutdown":
		return utils.Gray
	default:
		return nil
	}
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

var codespacesSuspendCmd = &cobra.Command{
	Use:   "suspend <codespacename>",
	Short: "Suspend a codespace",
	Example: heredoc.Doc(`
	$ gh codespaces suspend <codespacename>
	`),
	Args: cobra.ExactArgs(1),
	RunE: codespacesSuspend,
}

var codespacesResumeCmd = &cobra.Command{
	Use:   "resume <codespacename>",
	Short: "Resume a codespace",
	Example: heredoc.Doc(`
	$ gh codespaces resume <codespacename>
	`),
	Args: cobra.ExactArgs(1),
	RunE: codespacesResume,
}

var codespacesDeleteCmd = &cobra.Command{
	Use:   "delete <codespacename>",
	Short: "Delete a codespace",
	Example: heredoc.Doc(`
	$ gh codespaces delete <codespacename>
	`),
	Args: cobra.ExactArgs(1),
	RunE: codespacesDelete,
}

var codespacesCreateCmd = &cobra.Command{
	Use:   "create <repo>",
	Short: "Create a codespace",
	Example: heredoc.Doc(`
	$ gh codespaces create cli/cli
	`),
	Args: cobra.ExactArgs(1),
	RunE: codespacesCreate,
}
