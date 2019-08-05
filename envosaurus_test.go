package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tentwentyfive/envosaurus/specs"
)

func TestUnMarshallProject(t *testing.T) {
	git := specs.GitSpec{"foo"}
	var p = []specs.ProjectSpec{{"fred", nil}, {"wilma", &git}}
	var g []specs.ProjectSpec
	b := `[{"name": "fred"}, {"name": "wilma", "git": {"clone": "foo"}}]`
	json.Unmarshal([]byte(b), &g)

	if !reflect.DeepEqual(p, g) {
		t.Error("Expected ", p, "got ", g)
	}

}

func TestLoadProjects(t *testing.T) {
	path := "samples/projects.json"
	var projects specs.ProjectsSpec

	if err := projects.LoadProjects(path); err != nil {
		t.Error("Unable to load ", path, ": ", err)
	}

	git := specs.GitSpec{"git@github.com:kafkaex/kafka_ex"}
	expect := specs.ProjectsSpec{
		"/Users/dswain/envosrc",
		[]specs.ProjectSpec{
			{"KafkaEx", &git},
		},
	}

	if !reflect.DeepEqual(projects, expect) {
		t.Error("Expected ", expect, "got ", projects)
	}
}
