package cmd

import (
	"fmt"
	"os"

	"github.com/mucz/protobuf-decompiler/restore"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "protodec [file]",
	Short: "protodec help you recover protobuf files from binary descriptor in template code.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := restore.Do(args[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(data)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
