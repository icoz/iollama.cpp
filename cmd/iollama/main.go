package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "iollama",
		Short: "Lightweight OpenAI-compatible server for GGUF models",
	}

	// Команда serve (будет реализована в промпте 6)
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("serve command placeholder")
			return nil
		},
	}
	rootCmd.AddCommand(serveCmd)

	// Команда download (будет реализована в промпте 7)
	downloadCmd := &cobra.Command{
		Use:   "download [repo_id] [filename]",
		Short: "Download a GGUF model from Hugging Face Hub",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("download placeholder: %s %s\n", args[0], args[1])
			return nil
		},
	}
	rootCmd.AddCommand(downloadCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
