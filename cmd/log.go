/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		platform, err := cmd.Flags().GetString("platform")

		if err != nil {
			fmt.Println("You must provide a platform to watch logs", err)
			return
		}

		fmt.Println(platform)

		switch platform {
		case "android":
			watchAndroidLogs()
		}

	},
}

func watchAndroidLogs() {
	execCmd := exec.Command("adb", "logcat")

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe", err)
		return
	}

	if err := execCmd.Start(); err != nil {
		fmt.Println("Error starting command", err)
		return
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from StdoutPipe", err)
	}

	if err := execCmd.Wait(); err != nil {
		fmt.Println("Error wiating for command", err)
	}
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	logCmd.Flags().StringP("platform", "p", "", "Choose android or ios")
}
