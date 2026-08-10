package cli

import (
	"fmt"
	"sort"

	"github.com/UPin2905/portdoctor/pkg/framework"
	"github.com/UPin2905/portdoctor/pkg/output"
	"github.com/UPin2905/portdoctor/pkg/port"
	"github.com/UPin2905/portdoctor/pkg/process"
	"github.com/UPin2905/portdoctor/pkg/project"
	"github.com/UPin2905/portdoctor/pkg/runtime"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "List active listening ports on this machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Header()

		inspector := port.NewInspector()
		ports, err := inspector.ListListening()
		if err != nil {
			return fmt.Errorf("could not scan ports: %w", err)
		}

		if len(ports) == 0 {
			fmt.Println("No listening ports found.")
			return nil
		}

		// Sắp xếp theo port number
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Port < ports[j].Port
		})

		output.TableHeader()
		for _, p := range ports {
			procName := ""
			projName := ""
			typ := ""

			if p.PID > 0 {
				proc, err := process.GetProcess(p.PID)
				if err == nil && proc != nil {
					procName = proc.Name
					rt := runtime.Detect(proc.Name, proc.CommandLine)
					fw := framework.Detect(rt, proc.CommandLine, proc.WorkingDir)
					if fw != "" {
						typ = fw
					} else if rt != "" {
						typ = rt
					}
					proj := project.Detect(proc.WorkingDir)
					if proj != nil {
						projName = proj.Name
					}
				}
			}

			output.TableRow(p.Port, procName, p.PID, projName, typ)
		}
		fmt.Println()
		return nil
	},
	SilenceUsage: true,
}

