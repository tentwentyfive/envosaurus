package specs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

// ProjectsSpec describes a list of projects
type ProjectsSpec struct {
	RootDirectory string        `json:"rootDirectory"`
	Projects      []ProjectSpec `json:"projects"`
}

// GitSpec describes a git repo
type GitSpec struct {
	Clone string `json:"clone"`

	repo *git.Repository
}

// ProjectSpec describes a single project
type ProjectSpec struct {
	Name string  `json:"name"`
	Path string  `json:"path"`
	Git  GitSpec `json:"git,omitempty"`
}

// RepoFileIsReadable check if the file exists and is readable
func RepoFileIsReadable(path string) bool {
	_, err := os.Open(path)
	if err != nil {
		return false
	}
	return true
}

// LoadProjects load projects from the given path
func (projects *ProjectsSpec) LoadProjects(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteVal, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteVal, projects)

	return nil
}

// Write write projects to the given path
func (projects *ProjectsSpec) Write(path string) error {
	jsonData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	writeErr := ioutil.WriteFile(path, jsonData, 0644)
	return writeErr
}

// ContainsProjectAtPath returns true if the collection already contains
// a project at the same path
func (projects *ProjectsSpec) ContainsProjectAtPath(projectSpec *ProjectSpec) bool {
	for _, p := range projects.Projects {
		if p.Path == projectSpec.Path {
			return true
		}
	}
	return false
}

// RepoFromPath returns a ProjectSpec from a repo at the given path
func RepoFromPath(path string) (ProjectSpec, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		wrappedError := fmt.Errorf("Unable find a repository: %w", err)
		return ProjectSpec{Git: GitSpec{repo: repo}}, wrappedError
	}

	remotes, err := repo.Remotes()
	if err != nil {
		wrappedError := fmt.Errorf("Unable to determine remotes: %w", err)
		return ProjectSpec{Git: GitSpec{repo: repo}}, wrappedError
	}

	url := ""
	for _, remote := range remotes {
		config := remote.Config()
		if config.Name == "origin" {
			url = config.URLs[0]
		}
	}

	gitSpec := GitSpec{Clone: url, repo: repo}
	root, err := gitSpec.RepoRoot()
	if err != nil {
		return ProjectSpec{Git: gitSpec}, err
	}

	name := filepath.Base(root)
	return ProjectSpec{Name: name, Path: root, Git: gitSpec}, nil
}

// RepoRoot returns the root path of the repository
func (gitSpec *GitSpec) RepoRoot() (string, error) {
	// we already know we can get the worktree
	worktree, err := gitSpec.repo.Worktree()
	if err != nil {
		wrappedError := fmt.Errorf("Unable to determine repository path: %w", err)
		return "", wrappedError
	}

	return worktree.Filesystem.Root(), nil
}

// GetCloneOpts returns the path to clone into and git.CloneOptions values
func (projectSpec *ProjectSpec) GetCloneOpts(rootDir string) (string, git.CloneOptions, error) {
	toDir, err := filepath.Abs(filepath.Join(rootDir, projectSpec.Path))
	if err != nil {
		return toDir, git.CloneOptions{}, err
	}

	return toDir, git.CloneOptions{
		URL: projectSpec.Git.Clone,
	}, nil
}
