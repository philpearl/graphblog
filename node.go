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
func (nodes nodes) diameter() int {
	var diameter int
	q := &list{}
	bfsData := make([]bfsNode, len(nodes))
	for id := range nodes {
		df := nodes.longestShortestPath(nodeId(id), q, bfsData)
		if df > diameter {
			diameter = df
		}
		// Need to reset the bfsData between runs
		for i := range bfsData {
			d := &bfsData[i]
			d.depth = 0
			d.parent = nil
		}
	}
	return diameter
}

type bfsNode struct {
	// bfs tracking data
	parent *node
	depth  int
}

func (nodes nodes) longestShortestPath(start nodeId, q *list, bfsData []bfsNode) int {

	n := nodes.get(start)
	bfsData[n.id] = bfsNode{parent: n, depth: 0}
	q.pushBack(n)

	for {
		newN := q.getHead()
		if newN == nil {
			break
		}
		n = newN

		for _, id := range n.adj {
			bm := bfsData[id]
			if bm.parent == nil {
				bfsData[id] = bfsNode{parent: n, depth: bfsData[n.id].depth + 1}
				q.pushBack(nodes.get(id))
			}
		}
	}

	return bfsData[n.id].depth
}
