package codespaces

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/spf13/cobra"

	cmdCreate "github.com/cli/cli/pkg/cmd/codespaces/create"
	cmdDelete "github.com/cli/cli/pkg/cmd/codespaces/delete"
	cmdList "github.com/cli/cli/pkg/cmd/codespaces/list"
	cmdResume "github.com/cli/cli/pkg/cmd/codespaces/resume"
	cmdSuspend "github.com/cli/cli/pkg/cmd/codespaces/suspend"
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

	codespacesCmd.AddCommand(cmdList.NewCmdList(f))
	codespacesCmd.AddCommand(cmdSuspend.NewSuspendCmd(f))
	codespacesCmd.AddCommand(cmdResume.NewResumeCmd(f))
	codespacesCmd.AddCommand(cmdDelete.NewDeleteCmd(f))
	codespacesCmd.AddCommand(cmdCreate.NewCreateCmd(f))

	return codespacesCmd
}
