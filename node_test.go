package graphblog

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type edge struct{ a, b nodeName }
type edgeList []edge

func (e edgeList) build(g *graph) {
	for _, edge := range e {
		g.addEdge(edge.a, edge.b)
	}
}

func TestDiameter(t *testing.T) {

	tests := []struct {
		name        string
		edgeList    edgeList
		expDiameter int
	}{
		{
			name: "empty",
		},
		{
			name:        "1 edge",
			edgeList:    edgeList{{"a", "b"}},
			expDiameter: 1,
		},
		{
			name:        "3 in line",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}},
			expDiameter: 2,
		},
		{
			name:        "4 in line",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "d"}},
			expDiameter: 3,
		},
		{
			name:        "Triangle",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"a", "c"}},
			expDiameter: 1,
		},
		{
			name:        "Square",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "d"}, {"a", "d"}},
			expDiameter: 2,
		},
		{
			name:        "2 loops",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "a"}, {"c", "d"}, {"d", "e"}, {"e", "c"}},
			expDiameter: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := New(100)
			test.edgeList.build(g)
			g2 := g.convertToMemoryOptimisedGraph()
			dia := g2.diameter()
			if dia != test.expDiameter {
				t.Errorf("Diameter not as expected. Have %d, expected %d", dia, test.expDiameter)
			}
		})
	}
}

func BenchmarkDiameter(b *testing.B) {
	g := New(10000)
	// Load the test data
	f, err := os.Open("testdata/edges.txt")
	assert.NoError(b, err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		edge := strings.Fields(line)
		assert.Len(b, edge, 2)
		g.addEdge(nodeName(edge[0]), nodeName(edge[1]))
	}
	assert.NoError(b, err)

	g2 := g.convertToMemoryOptimisedGraph()

	b.Run("diameter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			d := g2.diameter()
			assert.Equal(b, 82, d)
		}
	})
}
