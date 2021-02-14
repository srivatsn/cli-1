package skus

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

type SkusOptions struct {
	HTTPClient  func() (*http.Client, error)
	IO          *iostreams.IOStreams
	ColorScheme *iostreams.ColorScheme
}

func NewSkusCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SkusOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	skusCmd := &cobra.Command{
		Use:   "skus",
		Short: "List all available skus",
		Example: heredoc.Doc(`
		$ gh codespaces skus
		`),
		Args: cmdutil.NoArgsQuoteReminder,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listSkus(opts)
		},
	}

	return skusCmd
}

func listSkus(opts *SkusOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	skus, err := api.ListSkus(apiClient, currentUser)
	if err != nil {
		return err
	}

	table := utils.NewTablePrinter(opts.IO)

	for _, sku := range skus {
		table.AddField(sku.Name, nil, opts.ColorScheme.Blue)
		table.AddField(sku.DisplayName, nil, opts.ColorScheme.Bold)
		table.AddField(sku.OperationSystem, nil, opts.ColorScheme.Cyan)
		table.EndRow()
	}

	err = table.Render()
	if err != nil {
		return err
	}

	return nil
}
