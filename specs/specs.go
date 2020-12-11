package specs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ProjectsSpec describes a list of projects
type ProjectsSpec struct {
	RootDirectory string        `json:"rootDirectory"`
	Projects      []ProjectSpec `json:"projects"`
}

// GitSpec describes a git repo
type GitSpec struct {
	Clone string `json:"clone"`
}

// ProjectSpec describes a single project
type ProjectSpec struct {
	Name string   `json:"name"`
	Git  *GitSpec `json:"git,omitempty"`
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

	projects.RootDirectory = os.ExpandEnv(projects.RootDirectory)
	return nil
}

func (projects *ProjectsSpec) Write(path string) error {
	jsonData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	writeErr := ioutil.WriteFile(path, jsonData, 0644)
	return writeErr
}

func (projects *ProjectsSpec) Contains(projectSpec *ProjectSpec) bool {
	for _, p := range projects.Projects {
		if p.Git.Clone == projectSpec.Git.Clone && p.Name == projectSpec.Name {
			return true
		}
	}
	return false
}
