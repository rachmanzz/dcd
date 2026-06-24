package data

import (
	"testing"
)

func TestDataSetGet(t *testing.T) {
	ds := NewDataSet(map[string]any{
		"info": map[string]any{
			"username": "john",
		},
	})
	val, ok := ds.Get("info.username")
	if !ok {
		t.Fatal("info.username not found")
	}
	if val != "john" {
		t.Fatalf("got %v, want john", val)
	}
}

func TestDataSetResolve(t *testing.T) {
	ds := NewDataSet(map[string]any{
		"info": map[string]any{
			"name": "John",
		},
	})
	out := ds.Resolve("hello {{info.name}}")
	if out != "hello John" {
		t.Fatalf("got %q, want hello John", out)
	}
}

func TestDataSetResolveMissing(t *testing.T) {
	ds := NewDataSet(map[string]any{})
	out := ds.Resolve("{{missing.key}}")
	if out != "{{missing.key}}" {
		t.Fatalf("got %q, want {{missing.key}}", out)
	}
}

func TestDataSetFromStruct(t *testing.T) {
	type Info struct {
		Username string
		Age      int
	}
	ds := NewDataSet(Info{Username: "john", Age: 25})
	val, ok := ds.Get("username")
	if !ok {
		t.Fatal("username not found")
	}
	if val != "john" {
		t.Fatalf("got %v, want john", val)
	}
}
