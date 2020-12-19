package specs

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestUnMarshallProject(t *testing.T) {
	git := GitSpec{Clone: "foo"}
	var p = []ProjectSpec{{Name: "fred", Git: GitSpec{}}, {Name: "wilma", Git: git}}
	var g []ProjectSpec
	b := `[{"name": "fred"}, {"name": "wilma", "git": {"clone": "foo"}}]`
	json.Unmarshal([]byte(b), &g)

	if !reflect.DeepEqual(p, g) {
		t.Error("Expected ", p, "got ", g)
	}

}

func TestLoadProjects(t *testing.T) {
	path := "../../../samples/projects.json"
	var projects ProjectsSpec

	if err := projects.LoadProjects(path); err != nil {
		t.Error("Unable to load ", path, ": ", err)
	}

	kafkaExGit := GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"}
	kayrockGit := GitSpec{Clone: "git@github.com:dantswain/kayrock"}
	kafkaExExamplesGit := GitSpec{Clone: "git@github.com:dantswain/kafka_ex_examples"}
	expect := ProjectsSpec{
		RootDirectory: "${HOME}/envosrc",
		Projects: []ProjectSpec{
			{Name: "KafkaEx", Git: kafkaExGit},
			{Name: "Kayrock", Git: kayrockGit},
			{Name: "KafkaExExamples", Git: kafkaExExamplesGit},
		},
	}

	if !reflect.DeepEqual(projects, expect) {
		t.Error("Expected ", expect, "got ", projects)
	}
}

func TestDetermineRepo(t *testing.T) {
	project, err := RepoFromPath(".")
	if err != nil {
		t.Error("Unexpected error ", err)
	}

	if project.Name != "envosaurus" {
		t.Error("Unexpected name ", project.Name)
	}

	if !strings.HasSuffix(project.Git.Clone, "tentwentyfive/envosaurus") {
		t.Error("Unexpected remote ", project.Git.Clone)
	}
}

func TestDetermineRepoInSubdirectory(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("Couldn't get working directory ", err)
	}

	project, err := RepoFromPath(path + "/internal/pkg/specs")
	if err != nil {
		t.Error("Unexpected error ", err)
	}

	if project.Name != "envosaurus" {
		t.Error("Unexpected name ", project.Name)
	}

	if !strings.HasSuffix(project.Git.Clone, "tentwentyfive/envosaurus") {
		t.Error("Unexpected remote ", project.Git.Clone)
	}
}
