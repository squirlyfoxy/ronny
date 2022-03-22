package main

import (
	"testing"

	dijkstra "github.com/squirlyfoxy/ronny/database/dijkstra"
)

func TestGraph(t *testing.T) {
	graph := dijkstra.NewGraph()
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 2)
	graph.AddEdge("b", "c", 3)
	graph.AddEdge("b", "d", 4)
	graph.AddEdge("c", "d", 5)
	graph.AddEdge("c", "e", 6)

	len, p := graph.GetPath("a", "b")
	if len != 1 {
		t.Error("GetPath() failed")
	}
	if p[0] != "a" {
		t.Error("GetPath() failed")
	}

	len, p = graph.GetPath("a", "c")
	if len != 2 {
		t.Error("GetPath() failed")
	}
	if p[0] != "a" || p[1] != "c" {
		t.Error("GetPath() failed")
	}
}
