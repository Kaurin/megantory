package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildHash should be injected during build time: -ldflags "-X main.buildHash=$(git rev-parse HEAD)"
var BuildHash string

// BuildVersion should be injected during build time: -ldflags "-X main.buildVersion=$(git describe --tags)"
var BuildVersion string

// BuildDate should be injected during build time: -ldflags "-X main.buildDate=$(date -u -Iseconds)"
var BuildDate string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints out the version information",
	Long:  `Prints out the version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", BuildVersion)
		fmt.Printf("Git Hash: %s\n", BuildHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
