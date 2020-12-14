package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tentwentyfive/envosaurus/internal/pkg/specs"
	"gopkg.in/src-d/go-git.v4"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a repository",
	Long:  `Add the repository in the current directory to your config`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.PlainOpen(".")
		if err != nil {
			fmt.Println("Error determining repo root")
			fmt.Println(err)
			os.Exit(1)
		}

		remotes, err := repo.Remotes()
		if err != nil {
			fmt.Println("Unable to list remotes")
			fmt.Println(err)
			os.Exit(1)
		}

		url := ""
		for _, remote := range remotes {
			config := remote.Config()
			if config.Name == "origin" {
				url = config.URLs[0]
			}
		}

		absPath, err := filepath.Abs(".")
		if err != nil {
			fmt.Println("Unable to determine absolute path")
			fmt.Println(err)
			os.Exit(1)
		}
		name := filepath.Base(absPath)

		gitSpec := specs.GitSpec{Clone: url}
		projectSpec := specs.ProjectSpec{Name: name, Git: &gitSpec}

		var projects specs.ProjectsSpec

		if specs.RepoFileIsReadable(repoConfig) {
			if err := projects.LoadProjects(repoConfig); err != nil {
				fmt.Println("Error loading projects")
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			// by default use the parent directory of the project being added
			projects.RootDirectory = filepath.Dir(absPath)
		}

		if projects.Contains(&projectSpec) {
			fmt.Println("Project already in file")
			os.Exit(1)
		}

		projects.Projects = append(projects.Projects, projectSpec)

		if err := projects.Write(repoConfig); err != nil {
			fmt.Println("Error writing config file")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s written to %s\n", name, repoConfig)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
