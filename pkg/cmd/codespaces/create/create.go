package create

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghinstance"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	HTTPClient  func() (*http.Client, error)
	IO          *iostreams.IOStreams
	ColorScheme *iostreams.ColorScheme
	RepoName    string
	RepoRef     string
	Sku         string
}

func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesCreateCmd := &cobra.Command{
		Use:   "create <repo>",
		Short: "Create a codespace",
		Example: heredoc.Doc(`
		$ gh codespaces create cli/cli
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RepoName = args[0]
			opts.RepoRef, _ = cmd.Flags().GetString("ref")
			opts.Sku, _ = cmd.Flags().GetString("sku")

			return createCodespace(opts, cmd, args)
		},
	}
	codespacesCreateCmd.Flags().StringP("ref", "r", "", "A ref in the repo from which the codespace will be created. The default branch of the repo is used otherwise")
	codespacesCreateCmd.Flags().StringP("sku", "s", "", "The sku of the codespace (eg: basicLinux)")

	return codespacesCreateCmd
}

func createCodespace(opts *CreateOptions, cmd *cobra.Command, args []string) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	codespaceName, err := api.CreateCodespace(apiClient, currentUser, opts.RepoName, opts.RepoRef, opts.Sku)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, opts.ColorScheme.Cyan("Codespace %s successfully created.\n"), codespaceName)
	return nil
}
