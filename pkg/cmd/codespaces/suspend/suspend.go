package suspend

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

type SuspendOptions struct {
	HTTPClient    func() (*http.Client, error)
	IO            *iostreams.IOStreams
	ColorScheme   *iostreams.ColorScheme
	CodespaceName string
}

func NewSuspendCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SuspendOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesSuspendCmd := &cobra.Command{
		Use:   "suspend <codespacename>",
		Short: "Suspend a codespace",
		Example: heredoc.Doc(`
		$ gh codespaces suspend <codespacename>
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CodespaceName = args[0]

			return suspendCodespace(opts)
		},
	}

	return codespacesSuspendCmd
}

func suspendCodespace(opts *SuspendOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	err = api.SuspendCodespace(apiClient, currentUser, opts.CodespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, opts.ColorScheme.Cyan("Codespace %s successfully suspended.\n"), opts.CodespaceName)
	return nil
}
