package cli

import (
	"fmt"
	"strconv"

	"github.com/UPin2905/portdoctor/pkg/framework"
	"github.com/UPin2905/portdoctor/pkg/output"
	"github.com/UPin2905/portdoctor/pkg/port"
	"github.com/UPin2905/portdoctor/pkg/process"
	"github.com/UPin2905/portdoctor/pkg/project"
	"github.com/UPin2905/portdoctor/pkg/runtime"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:     "inspect <port>",
	Short:   "Inspect what is using a port",
	Aliases: []string{"check"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInspect(args[0])
	},
	SilenceUsage: true,
}

func runInspect(portArg string) error {
	portNum, err := strconv.Atoi(portArg)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid port number", portArg)
	}
	if err := port.ValidatePort(portNum); err != nil {
		return err
	}

	output.Header()

	inspector := port.NewInspector()
	info, err := inspector.Inspect(portNum)
	if err != nil {
		return fmt.Errorf("could not inspect port %d: %w", portNum, err)
	}

	output.PortStatus(portNum, string(info.Status))

	if info.Status == port.StatusFree {
		fmt.Printf("Port %d is available.\n", portNum)
		return nil
	}

	if info.Status == port.StatusUnknown {
		fmt.Printf("Could not determine status of port %d.\n", portNum)
		return nil
	}

	// OCCUPIED: lấy process info
	if info.PID <= 0 {
		fmt.Println("Port is occupied but process information is unavailable.")
		output.SuggestedActions([]string{
			fmt.Sprintf("portdoctor kill %d", portNum),
			fmt.Sprintf("portdoctor find %d", portNum),
		})
		return nil
	}

	proc, err := process.GetProcess(info.PID)
	if err != nil || proc == nil {
		fmt.Printf("Port is occupied by PID %d (process details unavailable).\n", info.PID)
		return nil
	}

	// Process section
	output.Section("Process")
	output.Field("PID", strconv.Itoa(proc.PID))
	output.Field("Name", proc.Name)
	output.Field("Executable", proc.Executable)
	if !proc.StartTime.IsZero() {
		output.Field("Started", output.TimeAgo(proc.StartTime))
	}
	fmt.Println()

	// Command section
	if proc.CommandLine != "" {
		output.Section("Command")
		fmt.Printf("  %s\n\n", proc.CommandLine)
	}

	// Directory section
	if proc.WorkingDir != "" {
		output.Section("Directory")
		fmt.Printf("  %s\n\n", proc.WorkingDir)
	}

	// Parent section
	if proc.ParentPID > 0 {
		output.Section("Parent")
		parentStr := fmt.Sprintf("%s (PID %d)", proc.ParentName, proc.ParentPID)
		if proc.ParentName == "" {
			parentStr = fmt.Sprintf("PID %d", proc.ParentPID)
		}
		fmt.Printf("  %s\n\n", parentStr)
	}

	// Runtime + Framework detection
	rt := runtime.Detect(proc.Name, proc.CommandLine)
	fw := framework.Detect(rt, proc.CommandLine, proc.WorkingDir)

	if rt != "" || fw != "" {
		output.Section("Detected")
		output.Field("Runtime", rt)
		output.Field("Framework", fw)
		fmt.Println()
	}

	// Project detection
	proj := project.Detect(proc.WorkingDir)
	if proj != nil {
		output.Section("Project")
		output.Field("Name", proj.Name)
		output.Field("Directory", proj.Root)
		fmt.Println()
	}

	// Diagnosis
	diagnosis := buildDiagnosis(proc.Name, rt, fw, proj)
	if diagnosis != "" {
		output.Diagnosis(diagnosis)
	}

	output.SuggestedActions([]string{
		fmt.Sprintf("portdoctor kill %d", portNum),
		fmt.Sprintf("portdoctor find %d", portNum),
	})

	return nil
}

func buildDiagnosis(procName, rt, fw string, proj *project.ProjectInfo) string {
	switch {
	case fw != "" && proj != nil:
		return fmt.Sprintf("This appears to be a %s development server running from %s.", fw, proj.Root)
	case fw != "":
		return fmt.Sprintf("This appears to be a %s application.", fw)
	case rt != "":
		return fmt.Sprintf("This is a %s process (%s).", rt, procName)
	default:
		return fmt.Sprintf("Process: %s", procName)
	}
}

