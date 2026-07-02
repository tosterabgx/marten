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
		Use:          "tcp <port>",
		Short:        "Expose a local TCP port through the tunnel",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.ParseUint(args[0], 10, 16)
			if err != nil {
				return fmt.Errorf("invalid port %v", args[0])
			}

			return client.RunTunnel(uint16(port), false)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:          "http <port>",
		Short:        "Expose a local TCP port through the tunnel to http subdomain",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.ParseUint(args[0], 10, 16)
			if err != nil {
				return fmt.Errorf("invalid port %v", args[0])
			}

			return client.RunTunnel(uint16(port), true)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
