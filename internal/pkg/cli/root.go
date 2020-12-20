package cli

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envosaur",
	Short: "It helps you do stuff",
	Long:  `It helps you do a bunch of stuff`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	defaultRepoConfig := filepath.Join(usr.HomeDir, ".envosaurus", "repos.json")

	// common args
	rootCmd.PersistentFlags().StringVarP(&repoConfig, "repo-config", "r", defaultRepoConfig, "Path to repo config json file")
	rootCmd.MarkFlagRequired("repo-config")
}
