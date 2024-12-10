package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cyber_record_parser",
		Short: "A tool to parse cyber records",
	}

	var infoCmd = &cobra.Command{
		Use:   "info <record file>",
		Short: "Print the header information of the record file",
		Args:  cobra.ExactArgs(1),
		Run:   InfoCommand,
	}

    var topic string

	var echoCmd = &cobra.Command{
		Use:   "echo <record file> [--topic <topic>]",
		Short: "Print the messages of the record file",
		Args:  cobra.ExactArgs(1),
		Run:   EchoCommand,
	}

    echoCmd.Flags().StringVarP(&topic, "topic", "t", "", "The topic to print the messages")

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(echoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
