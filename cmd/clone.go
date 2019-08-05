package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tentwentyfive/envosaurus/specs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var repoConfig string

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories",
	Long:  `Clone repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		var projects specs.ProjectsSpec
		if err := projects.LoadProjects(repoConfig); err != nil {
			fmt.Println("Error loading projects")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Ensuring directory %s exists\n", projects.RootDirectory)
		os.MkdirAll(projects.RootDirectory, os.ModePerm)

		sshAuth, err := ssh.DefaultAuthBuilder("git")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, repo := range projects.Projects {
			parts := strings.Split(repo.Git.Clone, "/")
			lastPart := parts[len(parts)-1]

			toDir := filepath.Join(projects.RootDirectory, lastPart)

			fmt.Printf("\nCloning %s to %s\n", repo.Name, toDir)
			_, err := git.PlainClone(toDir, false, &git.CloneOptions{
				URL:      repo.Git.Clone,
				Progress: os.Stdout,
				Auth:     sshAuth,
			})

			if err == git.ErrRepositoryAlreadyExists {
				// this is ok, we just want to let the user know and continue
				fmt.Println(err)
			} else if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Println("\nAll repos cloned!")

	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	defaultRepoConfig := filepath.Join(usr.HomeDir, ".envosaurus", "repos.json")

	rootCmd.PersistentFlags().StringVarP(&repoConfig, "repo-config", "r", defaultRepoConfig, "Path to repo config json file")
	rootCmd.MarkFlagRequired("repo-config")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
