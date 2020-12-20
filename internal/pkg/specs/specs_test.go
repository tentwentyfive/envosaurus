package specs

import (
	"encoding/json"
	"os"
	"path/filepath"
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
			{Name: "KafkaEx", Path: "kafka/kafka_ex", Git: kafkaExGit},
			{Name: "Kayrock", Path: "kafka/kayrock", Git: kayrockGit},
			{Name: "KafkaExExamples", Path: "kafka/kafka_ex_examples", Git: kafkaExExamplesGit},
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

	fullPath, err := filepath.Abs("../../..")
	if err != nil {
		t.Error("Unexpected error ", err)
	}

	if project.Name != "envosaurus" {
		t.Error("Unexpected name ", project.Name)
	}

	if project.Path != fullPath {
		t.Error("Unexpected path ", project.Path, ", expected: ", fullPath)
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

	fullPath, err := filepath.Abs("../../..")
	if err != nil {
		t.Error("Unexpected error ", err)
	}

	project, err := RepoFromPath(path + "/internal/pkg/specs")
	if err != nil {
		t.Error("Unexpected error ", err)
	}

	if project.Name != "envosaurus" {
		t.Error("Unexpected name ", project.Name)
	}

	if project.Path != fullPath {
		t.Error("Unexpected path ", project.Path, ", expected: ", fullPath)
	}

	if !strings.HasSuffix(project.Git.Clone, "tentwentyfive/envosaurus") {
		t.Error("Unexpected remote ", project.Git.Clone)
	}
}

func TestGetCloneOptions(t *testing.T) {
	projectSpec := ProjectSpec{Path: "foo/subdir", Git: GitSpec{Clone: "git@github.com:foo/bar"}}

	toDir, opts, err := projectSpec.GetCloneOpts("/some/dir")
	if err != nil {
		t.Error("Unexpected error getting clone opts ", err)
	}

	if toDir != "/some/dir/foo/subdir" {
		t.Error("Unexpected toDir ", toDir)
	}

	if opts.URL != "git@github.com:foo/bar" {
		t.Error("Unexpected URL ", opts.URL)
	}
}

func TestSpecContainsProjectAtPath(t *testing.T) {
	kafkaExGit := GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"}
	kayrockGit := GitSpec{Clone: "git@github.com:dantswain/kayrock"}
	kafkaExExamplesGit := GitSpec{Clone: "git@github.com:dantswain/kafka_ex_examples"}
	projects := ProjectsSpec{
		RootDirectory: "${HOME}/envosrc",
		Projects: []ProjectSpec{
			{Name: "KafkaEx", Path: "kafka/kafka_ex", Git: kafkaExGit},
			{Name: "Kayrock", Path: "kafka/kayrock", Git: kayrockGit},
			{Name: "KafkaExExamples", Path: "kafka/kafka_ex_examples", Git: kafkaExExamplesGit},
		},
	}

	if !projects.ContainsProjectAtPath(&ProjectSpec{
		Name: "KafkaEx",
		Path: "kafka/kafka_ex",
		Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
	}) {
		t.Error("Project should have contained subproject at the same path")
	}

	if projects.ContainsProjectAtPath(&ProjectSpec{
		Name: "KafkaEx",
		Path: "other_dir/kafka_ex",
		Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
	}) {
		t.Error("Project should not have contained subproject at a different pat")
	}

	if !projects.ContainsProjectAtPath(&ProjectSpec{
		Name: "KafkaEx2",
		Path: "kafka/kafka_ex",
		Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
	}) {
		t.Error("Project should contain subproject with different name at the same path")
	}

	if !projects.ContainsProjectAtPath(&ProjectSpec{
		Name: "KafkaEx",
		Path: "kafka/kafka_ex",
		Git:  GitSpec{Clone: "git@github.com:dantswain/kafka_ex"},
	}) {
		t.Error("Project should contain subproject at the same path even if repo is different")
	}
}
