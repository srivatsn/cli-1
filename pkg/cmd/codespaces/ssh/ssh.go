package ssh

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghinstance"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type SSHOptions struct {
	HTTPClient    func() (*http.Client, error)
	IO            *iostreams.IOStreams
	ColorScheme   *iostreams.ColorScheme
	CodespaceName string
}

func NewSSHCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SSHOptions{
		HTTPClient:  f.HttpClient,
		IO:          f.IOStreams,
		ColorScheme: f.IOStreams.ColorScheme(),
	}

	codespacesSuspendCmd := &cobra.Command{
		Use:   "ssh <codespacename>",
		Short: "Connect to a codespace",
		Example: heredoc.Doc(`
		$ gh codespaces ssh <codespacename>
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CodespaceName = args[0]

			return sshToCodespace(opts)
		},
	}

	return codespacesSuspendCmd
}

func sshToCodespace(opts *SSHOptions) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	currentUser, err := api.CurrentLoginName(apiClient, ghinstance.OverridableDefault())
	if err != nil {
		return err
	}

	codespace, err := api.GetCodespaceDetails(apiClient, currentUser, opts.CodespaceName)
	if err != nil {
		return err
	}

	token, err := api.GetCodespaceToken(apiClient, currentUser, opts.CodespaceName)
	if err != nil {
		return err
	}
	fmt.Println(*token)

	workspaceID := codespace.Environment.Connection.SessionID
	fmt.Println(workspaceID)

	nodePath, err := exec.LookPath("node");
	if err != nil {
		return err
	}

	cmdNode := &exec.Cmd{
		Path: nodePath,
		Args: []string{nodePath, "app.js", "-w", workspaceID},
		Env: []string {"CODESPACE_TOKEN=" + *token},
		Stdout: opts.IO.Out,
		Stderr: opts.IO.ErrOut,
	}

	err = cmdNode.Start()
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second);
	
	sshConfig := &ssh.ClientConfig{
		User: "codespace",
		Auth: []ssh.AuthMethod{
			ssh.Password("testpwd1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", "localhost:2222", sshConfig)
	if err != nil {
		return err
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		return err
	}
	defer session.Close();

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	err = session.RequestPty("xterm", 40, 80, modes);
	if err != nil {
		return err;
	}

	session.Stdin = opts.IO.In;
	session.Stdout = opts.IO.Out;
	session.Stderr = opts.IO.ErrOut;

	// Start remote shell
	err = session.Shell();
	if err != nil {
		return err;
	}

	err = session.Wait();
	if err != nil {
		return err;
	}
	
	return nil
}