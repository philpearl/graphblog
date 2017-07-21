package graphblog

import (
	"runtime"
	"sync"
)

type nodeId int32
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
func (nl nodes) diameter() int {

	cpus := runtime.NumCPU()
	numNodes := len(nl)
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

				df := nl.longestShortestPath(nodeId(id), q, depths)
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

func (nodes nodes) longestShortestPath(start nodeId, q *list, depths []bfsNode) int {

	n := nodes.get(start)
	depths[n.id] = 0
	q.pushBack(n)

	for {
		newN := q.getHead()
		if newN == nil {
			break
		}
		n = newN

		for _, id := range n.adj {
			if depths[id] == -1 {
				depths[id] = depths[n.id] + 1
				q.pushBack(nodes.get(id))
			}
		}
	}

	return int(depths[n.id])
}
