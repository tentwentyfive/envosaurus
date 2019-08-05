package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnMarshallProject(t *testing.T) {
	git := GitSpec{"foo"}
	var p = []ProjectSpec{{"fred", nil}, {"wilma", &git}}
	var g []ProjectSpec
	b := `[{"name": "fred"}, {"name": "wilma", "git": {"clone": "foo"}}]`
	json.Unmarshal([]byte(b), &g)

	if !reflect.DeepEqual(p, g) {
		t.Error("Expected ", p, "got ", g)
	}

}

func TestLoadProjects(t *testing.T) {
	path := "samples/projects.json"
	var projects ProjectsSpec

	if err := projects.LoadProjects(path); err != nil {
		t.Error("Unable to load ", path, ": ", err)
	}

	git := GitSpec{"git@github.com:kafkaex/kafka_ex"}
	expect := ProjectsSpec{
		"${HOME}/envosrc",
		[]ProjectSpec{
			{"KafkaEx", &git},
		},
	}

	if !reflect.DeepEqual(projects, expect) {
		t.Error("Expected ", expect, "got ", projects)
	}
}
