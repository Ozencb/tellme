package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ResetApiKey bool
var Model string

var rootCmd = &cobra.Command{
	Use:   "tm",
	Short: "A CLI for Chat-GPT, written in Go",
	Args:  cobra.MinimumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&ResetApiKey, "reset", "r", false, "Reset API Key")
	rootCmd.PersistentFlags().StringVarP(&Model, "model", "m", "gpt-3.5-turbo", "GPT Model. Default: gpt-3.5-turbo")
}
