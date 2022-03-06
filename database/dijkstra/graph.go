package dijkstra

func NewGraph() *Graph {
	return &Graph{nodes: make(map[string][]Edge)}
}

func (g *Graph) AddEdge(origin, destiny string, weight int) {
	g.nodes[origin] = append(g.nodes[origin], Edge{node: destiny, weight: weight})
	g.nodes[destiny] = append(g.nodes[destiny], Edge{node: origin, weight: weight})
}

func (g *Graph) GetEdges(origin string) []Edge {
	return g.nodes[origin]
}

//Get the better path
func (g *Graph) GetPath(origin, destiny string) (int, []string) {
	h := NewHeap()
	h.Push(Path{value: 0, nodes: []string{origin}})
	visited := make(map[string]bool)

	for len(*h.values) > 0 {
		p := h.Pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == destiny {
			return p.value, p.nodes
		}

		for _, edge := range g.GetEdges(node) {
			if !visited[edge.node] {
				h.Push(Path{value: p.value + edge.weight, nodes: append([]string{}, append(p.nodes, edge.node)...)})
			}
		}

		visited[node] = true
	}

	return 0, nil
}
