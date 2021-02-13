package codespaces

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/spf13/cobra"

	codespacesListCmd "github.com/cli/cli/pkg/cmd/codespaces/list"
)

func NewCmdCodespaces(f *cmdutil.Factory) *cobra.Command {
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

	codespacesCmd.AddCommand(codespacesListCmd.NewCmdList(f))
	// codespacesCmd.AddCommand(codespacesSuspendCmd)
	// codespacesCmd.AddCommand(codespacesResumeCmd)
	// codespacesCmd.AddCommand(codespacesDeleteCmd)
	// codespacesCmd.AddCommand(codespacesCreateCmd)
	// codespacesCreateCmd.Flags().StringP("ref", "r", "", "A ref in the repo from which the codespace will be created. The default branch of the repo is used otherwise")
	// codespacesCreateCmd.Flags().StringP("sku", "s", "", "The sku of the codespace (eg: basicLinux)")

	return codespacesCmd
}

// func getAPIClientAndCurrentUser(f *cmdutil.Factory) (*api.Client, string, error) {
// 	httpClient, err := f.HttpClient()
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	apiClient := api.NewClientFromHTTP(httpClient)

// 	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return apiClient, currentUser, nil
// }

// func codespacesSuspend(cmd *cobra.Command, args []string) error {
// 	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
// 	if err != nil {
// 		return err
// 	}

// 	codespaceName := args[0]
// 	err = api.SuspendCodespace(apiClient, currentUser, codespaceName)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully suspended.\n"), codespaceName)
// 	return nil
// }

// func codespacesResume(cmd *cobra.Command, args []string) error {
// 	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
// 	if err != nil {
// 		return err
// 	}

// 	codespaceName := args[0]
// 	err = api.StartCodespace(apiClient, currentUser, codespaceName)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully resumed.\n"), codespaceName)
// 	return nil
// }

// func codespacesDelete(cmd *cobra.Command, args []string) error {
// 	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
// 	if err != nil {
// 		return err
// 	}

// 	codespaceName := args[0]
// 	err = api.DeleteCodespace(apiClient, currentUser, codespaceName)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully deleted.\n"), codespaceName)
// 	return nil
// }

// func codespacesCreate(cmd *cobra.Command, args []string) error {
// 	apiClient, currentUser, err := getAPIClientAndCurrentUser(cmd)
// 	if err != nil {
// 		return err
// 	}

// 	repoName := args[0]

// 	ref, err := cmd.Flags().GetString("ref")
// 	if err != nil {
// 		return err
// 	}

// 	sku, err := cmd.Flags().GetString("sku")
// 	if err != nil {
// 		return err
// 	}

// 	codespaceName, err := api.CreateCodespace(apiClient, currentUser, repoName, ref, sku)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(cmd.OutOrStdout(), utils.Cyan("Codespace %s successfully created.\n"), codespaceName)
// 	return nil
// }

// var codespacesSuspendCmd = &cobra.Command{
// 	Use:   "suspend <codespacename>",
// 	Short: "Suspend a codespace",
// 	Example: heredoc.Doc(`
// 	$ gh codespaces suspend <codespacename>
// 	`),
// 	Args: cobra.ExactArgs(1),
// 	RunE: codespacesSuspend,
// }

// var codespacesResumeCmd = &cobra.Command{
// 	Use:   "resume <codespacename>",
// 	Short: "Resume a codespace",
// 	Example: heredoc.Doc(`
// 	$ gh codespaces resume <codespacename>
// 	`),
// 	Args: cobra.ExactArgs(1),
// 	RunE: codespacesResume,
// }

// var codespacesDeleteCmd = &cobra.Command{
// 	Use:   "delete <codespacename>",
// 	Short: "Delete a codespace",
// 	Example: heredoc.Doc(`
// 	$ gh codespaces delete <codespacename>
// 	`),
// 	Args: cobra.ExactArgs(1),
// 	RunE: codespacesDelete,
// }

// var codespacesCreateCmd = &cobra.Command{
// 	Use:   "create <repo>",
// 	Short: "Create a codespace",
// 	Example: heredoc.Doc(`
// 	$ gh codespaces create cli/cli
// 	`),
// 	Args: cobra.ExactArgs(1),
// 	RunE: codespacesCreate,
// }
