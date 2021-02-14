package delete

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

type DeleteOptions struct {
	HTTPClient    func() (*http.Client, error)
	IO            *iostreams.IOStreams
	ColorScheme   *iostreams.ColorScheme
	CodespaceName string
}

func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesDeleteCmd := &cobra.Command{
		Use:   "delete <codespacename>",
		Short: "Delete a codespace",
		Example: heredoc.Doc(`
	$ gh codespaces delete <codespacename>
	`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CodespaceName = args[0]

			return deleteCodespace(opts)
		},
	}

	return codespacesDeleteCmd
}

func deleteCodespace(opts *DeleteOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	err = api.DeleteCodespace(apiClient, currentUser, opts.CodespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, opts.ColorScheme.Cyan("Codespace %s successfully deleted.\n"), opts.CodespaceName)
	return nil
}
