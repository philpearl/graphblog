package graphblog

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
	adj map[nodeId]*node
}

func (n *node) add(adjNode *node) {
	n.adj[adjNode.id] = adjNode
}

type nodes []node

func (nl nodes) init() {
	for i := range nl {
		nl[i].id = nodeId(i)
		nl[i].adj = make(map[nodeId]*node)
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
func (nodes nodes) diameter() int {
	var diameter int
	for id := range nodes {
		df := nodes.longestShortestPath(nodeId(id))
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

	bfsData := make([]bfsNode, len(nodes))

	n := nodes.get(start)
	bfsData[n.id] = bfsNode{parent: n, depth: 0}
	q.pushBack(n)

	for {
		newN := q.getHead()
		if newN == nil {
			break
		}
		n = newN

		for id, m := range n.adj {
			bm := bfsData[id]
			if bm.parent == nil {
				bfsData[id] = bfsNode{parent: n, depth: bfsData[n.id].depth + 1}
				q.pushBack(m)
			}
		}
	}

	return bfsData[n.id].depth
}
