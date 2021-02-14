package resume

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

type ResumeOptions struct {
	HTTPClient    func() (*http.Client, error)
	IO            *iostreams.IOStreams
	ColorScheme   *iostreams.ColorScheme
	CodespaceName string
}

func NewResumeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ResumeOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesResumeCmd := &cobra.Command{
		Use:   "resume <codespacename>",
		Short: "Resume a codespace",
		Example: heredoc.Doc(`
		$ gh codespaces resume <codespacename>
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CodespaceName = args[0]

			return resumeCodespace(opts)
		},
	}

	return codespacesResumeCmd
}

func resumeCodespace(opts *ResumeOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	err = api.StartCodespace(apiClient, currentUser, opts.CodespaceName)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, opts.ColorScheme.Cyan("Codespace %s successfully resumed.\n"), opts.CodespaceName)
	return nil
}
