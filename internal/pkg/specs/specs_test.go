package specs

import (
	"encoding/json"
	"log"
	"os/user"
	"reflect"
	"testing"
)

func TestUnMarshallProject(t *testing.T) {
	git := GitSpec{Clone: "foo"}
	var p = []ProjectSpec{{Name: "fred", Git: nil}, {Name: "wilma", Git: &git}}
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

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if err := projects.LoadProjects(path); err != nil {
		t.Error("Unable to load ", path, ": ", err)
	}

	kafkaExGit := GitSpec{Clone: "git@github.com:kafkaex/kafka_ex"}
	kayrockGit := GitSpec{Clone: "git@github.com:dantswain/kayrock"}
	kafkaExExamplesGit := GitSpec{Clone: "git@github.com:dantswain/kafka_ex_examples"}
	expect := ProjectsSpec{
		RootDirectory: usr.HomeDir + "/envosrc",
		Projects: []ProjectSpec{
			{Name: "KafkaEx", Git: &kafkaExGit},
			{Name: "Kayrock", Git: &kayrockGit},
			{Name: "KafkaExExamples", Git: &kafkaExExamplesGit},
		},
	}

	if !reflect.DeepEqual(projects, expect) {
		t.Error("Expected ", expect, "got ", projects)
	}
}
