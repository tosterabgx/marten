package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tosterabgx/marten/internal/client"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "marten",
		Short: "marten is a reverse TCP tunnel",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "tcp <port>",
		Short: "Expose a local TCP port through the tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.ParseUint(args[0], 10, 16)
			if err != nil {
				return fmt.Errorf("invalid port %v", args[0])
			}

			return client.RunTCPTunnel(uint16(port))
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
