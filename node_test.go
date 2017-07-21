package graphblog

import (
	"os"
	"testing"

	"bufio"

	"strings"

	"github.com/stretchr/testify/assert"
)

type edge struct{ a, b nodeId }
type graph []edge

func (g graph) build(nodes nodes) {
	for _, edge := range g {
		nodes.addEdge(edge.a, edge.b)
	}
}

func TestDiameter(t *testing.T) {

	tests := []struct {
		name        string
		graph       graph
		expDiameter int
	}{
		{
			name: "empty",
		},
		{
			name:        "1 edge",
			graph:       graph{{"a", "b"}},
			expDiameter: 1,
		},
		{
			name:        "3 in line",
			graph:       graph{{"a", "b"}, {"b", "c"}},
			expDiameter: 2,
		},
		{
			name:        "4 in line",
			graph:       graph{{"a", "b"}, {"b", "c"}, {"c", "d"}},
			expDiameter: 3,
		},
		{
			name:        "Triangle",
			graph:       graph{{"a", "b"}, {"b", "c"}, {"a", "c"}},
			expDiameter: 1,
		},
		{
			name:        "Square",
			graph:       graph{{"a", "b"}, {"b", "c"}, {"c", "d"}, {"a", "d"}},
			expDiameter: 2,
		},
		{
			name:        "2 loops",
			graph:       graph{{"a", "b"}, {"b", "c"}, {"c", "a"}, {"c", "d"}, {"d", "e"}, {"e", "c"}},
			expDiameter: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes := make(nodes)
			test.graph.build(nodes)
			dia := nodes.diameter()
			if dia != test.expDiameter {
				t.Errorf("Diameter not as expected. Have %d, expected %d", dia, test.expDiameter)
			}
		})
	}
}

func BenchmarkDiameter(b *testing.B) {
	nodes := make(nodes)
	// Load the test data
	f, err := os.Open("testdata/edges.txt")
	assert.NoError(b, err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		edge := strings.Fields(line)
		assert.Len(b, edge, 2)
		nodes.addEdge(nodeId(edge[0]), nodeId(edge[1]))
	}
	assert.NoError(b, err)

	b.Run("diameter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			d := nodes.diameter()
			assert.Equal(b, 82, d)
		}
	})
}
