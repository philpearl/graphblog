package graphblog

import (
	"runtime"
	"sync"
)

type nodeId int
type nodeName string

type symbolTable map[nodeName]nodeId

func (s symbolTable) getId(name nodeName) nodeId {
	id, ok := s[name]
	if !ok {
		id = nodeId(len(s))
		s[name] = id
	}
	return id
}

type graph struct {
	symbolTable
	nodes
	noEdges int32
}

type graph2 struct {
	// The node ID is the index into the eStartIdx array
	eStartIdx []nodeId // array size is no of nodes +1
	edges     []nodeId // array size is no of edges *2 (as each edge is forward/backward)
}

func (g *graph) convertToMemoryOptimisedGraph() *graph2 {
	numNodes := len(g.nodes)
	numEdges := g.noEdges

	g2 := &graph2{}
	g2.eStartIdx = make([]nodeId, numNodes+1, numNodes+1)
	g2.edges = make([]nodeId, numEdges*2, numEdges*2)

	var edgesPos nodeId
	for id := 0; id < numNodes; id++ {
		g2.eStartIdx[id] = edgesPos
		node := g.nodes[id]
		for _, adjID := range node.adj {
			g2.edges[edgesPos] = adjID
			edgesPos++
		}
	}
	g2.eStartIdx[numNodes] = edgesPos

	return g2.reOrder()
}

func (g2 *graph2) traversalOrder() []nodeId {
	numNodes := len(g2.eStartIdx) - 1
	if numNodes == 0 {
		return nil
	}

	q := &list{}
	visited := make([]bool, numNodes)

	var n nodeId // starting node 0 so no need to initalise it

	order := make([]nodeId, 0, numNodes)

	for n = 0; n < nodeId(numNodes); n++ { // Loop through all component ids -
		if visited[n] {
			continue
		}
		q.pushBack(n)
		visited[n] = true
		for { // Do BFS

			nodeID := q.getHead()
			if nodeID == -1 {
				break
			}
			order = append(order, nodeID)
			for edgesIdx := g2.eStartIdx[nodeID]; edgesIdx < g2.eStartIdx[nodeID+1]; edgesIdx++ {
				id := g2.edges[edgesIdx]
				if !visited[id] {
					visited[id] = true
					q.pushBack(id)
				}
			}
		}
	}
	return order
}

func (g2 *graph2) reOrder() *graph2 {
	numNodes := len(g2.eStartIdx) - 1

	g2New := &graph2{}
	g2New.eStartIdx = make([]nodeId, len(g2.eStartIdx), len(g2.eStartIdx))
	g2New.edges = make([]nodeId, len(g2.edges), len(g2.edges))

	// Get the order we traverse the graph in and then use that to reorder the nodeID and edges
	// to that ordering - idea is that it is more cache efficent.
	order := g2.traversalOrder()

	// 1st node id is the same but then we rebuid
	mapOldNew := make(map[nodeId]nodeId)
	for newID, oldID := range order {
		mapOldNew[oldID] = nodeId(newID)
	}
	var edgesPos nodeId
	for newID, oldID := range order {
		g2New.eStartIdx[newID] = edgesPos
		edgesStartIdx := g2.eStartIdx[oldID]
		edgesEndIdx := g2.eStartIdx[oldID+1]
		for edgesIdx := edgesStartIdx; edgesIdx < edgesEndIdx; edgesIdx++ {
			oldID := g2.edges[edgesIdx]
			g2New.edges[edgesPos] = mapOldNew[oldID]
			edgesPos++
		}
	}
	g2New.eStartIdx[numNodes] = edgesPos

	return g2New
}

func New(numNodes int) *graph {
	g := &graph{
		symbolTable: make(symbolTable, numNodes),
		nodes:       make(nodes, numNodes),
	}
	g.nodes.init()
	return g
}

func (g *graph) addEdge(a, b nodeName) {
	aid := g.symbolTable.getId(a)
	bid := g.symbolTable.getId(b)
	g.noEdges++
	g.nodes.addEdge(aid, bid)
}

type node struct {
	id nodeId

	// adjacent edges
	adj []nodeId
}

func (n *node) add(adjNode *node) {
	for _, id := range n.adj {
		if id == adjNode.id {
			return
		}
	}
	n.adj = append(n.adj, adjNode.id)
}

type nodes []node

func (nl nodes) init() {
	for i := range nl {
		nl[i].id = nodeId(i)
	}
}

func (nl nodes) get(id nodeId) *node {
	return &nl[id]
}

func (nl nodes) addEdge(a, b nodeId) {
	an := nl.get(a)
	bn := nl.get(b)

	an.add(bn)
	bn.add(an)
}

// diameter is the maximum length of a shortest path in the network
func (g2 *graph2) diameter() int {

	cpus := runtime.NumCPU()
	numNodes := len(g2.eStartIdx) - 1
	nodesPerCpu := numNodes / cpus

	results := make([]int, cpus)
	wg := &sync.WaitGroup{}
	wg.Add(cpus)
	start := 0
	for cpu := 0; cpu < cpus; cpu++ {
		end := start + nodesPerCpu
		if cpu == cpus-1 {
			end = numNodes
		}

		go func(cpu int, start, end nodeId) {
			defer wg.Done()
			var diameter int
			q := &list{}
			depths := make([]bfsNode, numNodes)
			for id := start; id < end; id++ {
				// Need to reset the bfsData between runs
				for i := range depths {
					depths[i] = -1
				}

				df := g2.longestShortestPath(nodeId(id), q, depths)
				if df > diameter {
					diameter = df
				}
			}
			results[cpu] = diameter
		}(cpu, nodeId(start), nodeId(end))
		start += nodesPerCpu
	}

	wg.Wait()

	diameter := 0
	for _, result := range results {
		if result > diameter {
			diameter = result
		}
	}
	return diameter
}

// bfs tracking data
type bfsNode int16

func (g2 *graph2) longestShortestPath(start nodeId, q *list, depths []bfsNode) int {

	n := start
	depths[n] = 0
	q.pushBack(n)

	for {
		newN := q.getHead()
		if newN < 0 {
			break
		}
		n = newN

		edgesStartIdx := g2.eStartIdx[n]
		edgesEndIdx := g2.eStartIdx[n+1]
		for edgesIdx := edgesStartIdx; edgesIdx < edgesEndIdx; edgesIdx++ {

			id := g2.edges[edgesIdx]

			if depths[id] == -1 {
				depths[id] = depths[n] + 1
				q.pushBack(id)
			}
		}
	}

	return int(depths[n])
}
