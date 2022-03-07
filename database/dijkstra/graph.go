//From: https://dev.to/douglasmakey/implementation-of-dijkstra-using-heap-in-go-6e3

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

	if origin == destiny {
		return 0, []string{origin}
	}

	if len(g.nodes[origin]) == 0 {
		return -1, []string{"origin"}
	}

	if len(g.nodes[destiny]) == 0 {
		return -1, []string{"destiny"}
	}

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

	return -1, nil
}
