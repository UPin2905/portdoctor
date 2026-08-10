package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/UPin2905/portdoctor/internal/port"
	"github.com/UPin2905/portdoctor/internal/process"
	"github.com/spf13/cobra"
)

var forceKill bool

var killCmd = &cobra.Command{
	Use:   "kill <port>",
	Short: "Terminate the process using a port (with confirmation)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("'%s' is not a valid port number", args[0])
		}
		if err := port.ValidatePort(portNum); err != nil {
			return err
		}

		inspector := port.NewInspector()
		info, err := inspector.Inspect(portNum)
		if err != nil {
			return fmt.Errorf("could not inspect port %d: %w", portNum, err)
		}

		if info.Status == port.StatusFree {
			fmt.Printf("Port %d is not in use.\n", portNum)
			return nil
		}

		if info.PID <= 0 {
			return fmt.Errorf("could not determine which process is using port %d", portNum)
		}

		proc, err := process.GetProcess(info.PID)
		if err != nil || proc == nil {
			return fmt.Errorf("could not get process info for PID %d: %w", info.PID, err)
		}

		// Hiển thị thông tin process trước khi hỏi
		fmt.Printf("Port %d is used by:\n\n", portNum)
		fmt.Printf("  %s\n", proc.Name)
		fmt.Printf("  PID %d\n", proc.PID)
		if proc.WorkingDir != "" {
			fmt.Printf("  %s\n", proc.WorkingDir)
		}
		fmt.Println()

		if !forceKill {
			fmt.Print("Terminate this process? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := killProcess(proc.PID); err != nil {
			return fmt.Errorf("could not terminate PID %d: %w", proc.PID, err)
		}

		fmt.Printf("Process %d (%s) has been terminated.\n", proc.PID, proc.Name)
		return nil
	},
	SilenceUsage: true,
}

func init() {
	killCmd.Flags().BoolVar(&forceKill, "force", false, "Skip confirmation prompt")
}

func killProcess(pid int) error {
	return killPlatform(pid)
}

// killPlatform adalah platform-agnostic wrapper.
// Implementation detail: pada Windows, os.FindProcess + Kill mengirim TerminateProcess.
// Pada Unix, mengirim SIGTERM dulu.
func killPlatform(pid int) error {
	if runtime.GOOS == "windows" {
		return killWindows(pid)
	}
	return killUnix(pid)
}

