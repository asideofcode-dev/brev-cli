// notebook/notebook.go
package notebook

import (
	"github.com/brevdev/brev-cli/pkg/cmd/hello"
	"github.com/brevdev/brev-cli/pkg/cmd/portforward"
	"github.com/brevdev/brev-cli/pkg/cmd/util"
	"github.com/brevdev/brev-cli/pkg/entity"
	breverrors "github.com/brevdev/brev-cli/pkg/errors"
	"github.com/brevdev/brev-cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	notebookLong    = "Open a notebook on your Brev machine"
	notebookExample = "brev notebook <InstanceName>"
)

type NotebookStore interface {
	portforward.PortforwardStore
}

type WorkspaceResult struct {
	Workspace *entity.Workspace // Replace with the actual type of workspace returned by GetUserWorkspaceByNameOrIDErr
	Err       error
}

func NewCmdNotebook(store NotebookStore, _ *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notebook",
		Short:   "Open a notebook on your Brev machine",
		Long:    notebookLong,
		Example: notebookExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Channel to get the result of the network call
			resultCh := make(chan *WorkspaceResult)

			// Start the network call in a goroutine
			go func() {
				workspace, err := util.GetUserWorkspaceByNameOrIDErr(store, args[0])
				resultCh <- &WorkspaceResult{Workspace: workspace, Err: err}
			}()

			// Type out the checking message
			hello.TypeItToMeUnskippable27("Checking to make sure the workspace is running...")

			// Wait for the network call to finish
			result := <-resultCh

			if result.Err != nil {
				return breverrors.WrapAndTrace(result.Err)
			}

			// Check if the workspace is running
			if result.Workspace.Status != "RUNNING" {
				hello.TypeItToMeUnskippable27("The workspace is not running. Please ensure it's in the running state before proceeding.")
				return breverrors.WorkspaceNotRunning{Status: result.Workspace.Status}
			}

			urlType := color.New(color.FgCyan, color.Bold).SprintFunc()
			warningType := color.New(color.FgBlack, color.Bold, color.BgCyan).SprintFunc()

			hello.TypeItToMeUnskippable27("\n" + warningType("  Please keep this terminal open 🤙  "))

			hello.TypeItToMeUnskippable27("\nClick here to go to your Jupyter notebook:\n\t 👉" + urlType("http://localhost:8888") + "👈\n\n\n")

			// Port forward on 8888
			err2 := portforward.RunPortforward(store, args[0], "8888:8888")
			if err2 != nil {
				return breverrors.WrapAndTrace(err2)
			}

			// Print out a link for the user
			hello.TypeItToMeUnskippable27("Your notebook is accessible at: http://localhost:8888")

			return nil
		},
	}

	return cmd
}
