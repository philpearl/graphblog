package graphblog

import (
	"container/list"
)

type nodeId string

type node struct {
	id nodeId

	// adjacent edges
	adj map[nodeId]*node
}

func (n *node) add(adjNode *node) {
	n.adj[adjNode.id] = adjNode
}

type nodes map[nodeId]*node

func (nodes nodes) get(id nodeId) *node {
	n, ok := nodes[id]
	if !ok {
		n = &node{
			id:  id,
			adj: make(map[nodeId]*node),
		}
		nodes[id] = n
	}
	return n
}

func (nodes *nodes) addEdge(a, b nodeId) {
	an := nodes.get(a)
	bn := nodes.get(b)

	an.add(bn)
	bn.add(an)
}

// diameter is the maximum length of a shortest path in the network
func (nodes nodes) diameter() int {
	var diameter int
	for id := range nodes {
		df := nodes.longestShortestPath(id)
		if df > diameter {
			diameter = df
		}
	}
	return diameter
}

type bfsNode struct {
	// bfs tracking data
	parent *node
	depth  int
}

func (nodes nodes) longestShortestPath(start nodeId) int {
	q := list.New()

	bfsData := make(map[nodeId]bfsNode, len(nodes))

	n := nodes.get(start)
	bfsData[n.id] = bfsNode{parent: n, depth: 0}
	q.PushBack(n)

	for {
		elt := q.Front()
		if elt == nil {
			break
		}
		n = q.Remove(elt).(*node)

		for id, m := range n.adj {
			bm := bfsData[id]
			if bm.parent == nil {
				bfsData[id] = bfsNode{parent: n, depth: bfsData[n.id].depth + 1}
				q.PushBack(m)
			}
		}
	}

	return bfsData[n.id].depth
}
