package cli

import (
	"fmt"
	"net"
	"strconv"

	"github.com/portdoctor/portdoctor/internal/output"
	"github.com/portdoctor/portdoctor/internal/port"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <port>",
	Short: "Find the nearest available port starting from the given port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("'%s' is not a valid port number", args[0])
		}
		if err := port.ValidatePort(portNum); err != nil {
			return err
		}

		output.Header()

		free, err := findFreePort(portNum)
		if err != nil {
			return err
		}

		if free == portNum {
			fmt.Printf("Port %d is available.\n", portNum)
		} else {
			fmt.Printf("Port %d is occupied.\n\n", portNum)
			fmt.Printf("Nearest available port:\n  %s\n", output.Highlight(strconv.Itoa(free)))
		}
		return nil
	},
	SilenceUsage: true,
}

// findFreePort tìm port trống gần nhất bắt đầu từ start.
func findFreePort(start int) (int, error) {
	for p := start; p <= 65535 && p < start+100; p++ {
		if isPortFree(p) {
			return p, nil
		}
	}
	return 0, fmt.Errorf("no free port found in range %d–%d", start, start+100)
}

func isPortFree(p int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
