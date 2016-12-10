package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var logRequests bool
var host string
var httpPort int

const defaultHttpPort = 8081
const defaultReqLog = true

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rest-api-skeleton",
	Short: "The rest api skeleton application",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
}
