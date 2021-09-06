package importgraph_test

import (
	"testing"

	"github.com/200sc/bebop/internal/importgraph"
)

func TestDgraphCycle(t *testing.T) {
	t.Run("cycle", func(t *testing.T) {
		dg := importgraph.NewDgraph()
		dg.AddEdge("a", "b")
		dg.AddEdge("b", "c")
		dg.AddEdge("c", "d")
		dg.AddEdge("d", "a")
		err := dg.FindCycle()
		if err == nil {
			t.Fatal("expected import cycle, got no error")
		}
	})
	t.Run("no cycle", func(t *testing.T) {
		dg := importgraph.NewDgraph()
		dg.AddEdge("a", "b")
		dg.AddEdge("b", "c")
		dg.AddEdge("c", "d")
		dg.AddEdge("e", "a")
		err := dg.FindCycle()
		if err != nil {
			t.Fatal("expected no import cycle, got error")
		}
	})
}
