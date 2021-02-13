package list

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghinstance"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	HTTPClient  func() (*http.Client, error)
	IO          *iostreams.IOStreams
	ColorScheme *iostreams.ColorScheme
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesListCmd := &cobra.Command{
		Use:   "list",
		Short: "List codespaces for current user",
		Args:  cmdutil.NoArgsQuoteReminder,
		Example: heredoc.Doc(`
		$ gh codespaces list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCodespaces(opts)
		},
	}

	return codespacesListCmd
}

func listCodespaces(opts *ListOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	response, err := api.GetCodespaces(apiClient, currentUser)
	if err != nil {
		return err
	}

	table := utils.NewTablePrinter(opts.IO)

	for _, codespace := range response.Codespaces {
		codespaceDetails, err := api.GetCodespaceDetails(apiClient, currentUser, codespace.Name)
		if err != nil {
			return err
		}

		hasUnpushedChanges := ""
		if codespaceDetails.Environment.HasUnpushedChanges {
			hasUnpushedChanges = "Has unpushed changes"
		}

		table.AddField(codespace.Name, nil, colorfuncForState(opts.ColorScheme, codespaceDetails.Environment.State))
		table.AddField(codespaceDetails.Environment.State, nil, colorfuncForState(opts.ColorScheme, codespaceDetails.Environment.State))
		table.AddField(codespaceDetails.Environment.SkuDisplayName, nil, nil)
		table.AddField(codespaceDetails.Environment.Seed.Moniker, nil, opts.ColorScheme.Blue)
		table.AddField(hasUnpushedChanges, nil, opts.ColorScheme.Red)
		table.EndRow()
	}

	err = table.Render()
	if err != nil {
		return err
	}

	return nil
}

func colorfuncForState(cs *iostreams.ColorScheme, state string) func(string) string {
	switch state {
	case "Available":
		return cs.Green
	case "Shutdown":
		return cs.Gray
	default:
		return nil
	}
}
