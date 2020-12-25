package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tentwentyfive/envosaurus/internal/pkg/specs"
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

		if !specs.RepoFileIsReadable(repoConfig) {
			fmt.Printf("Unable to load config file %s\n", repoConfig)
			os.Exit(1)
		}

		if err := projects.LoadProjects(repoConfig); err != nil {
			fmt.Println("Error loading projects")
			fmt.Println(err)
			os.Exit(1)
		}

		rootDir := os.ExpandEnv(projects.RootDirectory)
		fmt.Printf("Ensuring directory %s exists\n", rootDir)
		os.MkdirAll(rootDir, os.ModePerm)

		sshAuth, err := ssh.DefaultAuthBuilder("git")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, repo := range projects.Projects {
			toDir, cloneOpts, err := repo.GetCloneOpts(rootDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			cloneOpts.Progress = os.Stdout
			cloneOpts.Auth = sshAuth

			fmt.Printf("\nCloning %s to %s\n", repo.Name, toDir)
			_, err = git.PlainClone(toDir, false, &cloneOpts)

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
}
