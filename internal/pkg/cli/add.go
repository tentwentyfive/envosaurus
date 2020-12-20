package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tentwentyfive/envosaurus/internal/pkg/specs"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a repository",
	Long:  `Add the repository in the current directory to your config`,
	Run: func(cmd *cobra.Command, args []string) {
		projectSpec, err := specs.RepoFromPath(".")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		root, err := projectSpec.Git.RepoRoot()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var projects specs.ProjectsSpec

		if specs.RepoFileIsReadable(repoConfig) {
			if err := projects.LoadProjects(repoConfig); err != nil {
				fmt.Println("Error loading projects")
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			// by default use the parent directory of the project being added
			projects.RootDirectory = filepath.Dir(root)
		}

		relPath, err := filepath.Rel(projects.RootDirectory, projectSpec.Path)
		if err != nil {
			fmt.Println("Unable to compute relative path")
			fmt.Println(err)
			os.Exit(1)
		}
		projectSpec.Path = relPath

		if projects.ContainsProjectAtPath(&projectSpec) {
			fmt.Println("The project file already contains a project at the path ", projectSpec.Path)
			os.Exit(1)
		}

		projects.Projects = append(projects.Projects, projectSpec)

		if err := projects.Write(repoConfig); err != nil {
			fmt.Println("Error writing config file")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s written to %s\n", projectSpec.Name, repoConfig)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
