package specs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnMarshallProject(t *testing.T) {
	git := GitSpec{Clone: "foo"}
	var p = []ProjectSpec{{Name: "fred", Git: GitSpec{}}, {Name: "wilma", Git: git}}
	var g []ProjectSpec
	b := `[{"name": "fred"}, {"name": "wilma", "git": {"clone": "foo"}}]`
	json.Unmarshal([]byte(b), &g)

	assert.Equal(t, p, g)
}

func TestLoadProjects(t *testing.T) {
	path := "../../../samples/projects.json"
	var projects ProjectsSpec

	err := projects.LoadProjects(path)
	require.NoError(t, err)

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

	assert.Equal(t, expect, projects)
}

func TestDetermineRepo(t *testing.T) {
	project, err := RepoFromPath(".")
	require.NoError(t, err)

	fullPath, err := filepath.Abs("../../..")
	require.NoError(t, err)

	assert.Equal(t, "envosaurus", project.Name)
	assert.Equal(t, fullPath, project.Path)
	assert.Regexp(t, regexp.MustCompile("tentwentyfive/envosaurus$"), project.Git.Clone)
}

func TestDetermineRepoInSubdirectory(t *testing.T) {
	path, err := os.Getwd()
	require.NoError(t, err)

	fullPath, err := filepath.Abs("../../..")
	require.NoError(t, err)

	project, err := RepoFromPath(path + "/internal/pkg/specs")
	require.NoError(t, err)

	assert.Equal(t, "envosaurus", project.Name)
	assert.Equal(t, fullPath, project.Path)
	assert.Regexp(t, regexp.MustCompile("tentwentyfive/envosaurus$"), project.Git.Clone)
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

	assert.True(
		t,
		projects.ContainsProjectAtPath(&ProjectSpec{
			Name: "KafkaEx",
			Path: "kafka/kafka_ex",
			Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
		}), "Project should have contained subproject at the same path")

	assert.False(
		t,
		projects.ContainsProjectAtPath(&ProjectSpec{
			Name: "KafkaEx",
			Path: "other_dir/kafka_ex",
			Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
		}), "Project should not have contained subproject at a different pat")

	assert.True(
		t,
		projects.ContainsProjectAtPath(&ProjectSpec{
			Name: "KafkaEx2",
			Path: "kafka/kafka_ex",
			Git:  GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"},
		}), "Project should contain subproject with different name at the same path")

	assert.True(
		t,
		projects.ContainsProjectAtPath(&ProjectSpec{
			Name: "KafkaEx",
			Path: "kafka/kafka_ex",
			Git:  GitSpec{Clone: "git@github.com:dantswain/kafka_ex"},
		}), "Project should contain subproject at the same path even if repo is different")
}
