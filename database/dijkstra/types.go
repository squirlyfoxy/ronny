//From: https://dev.to/douglasmakey/implementation-of-dijkstra-using-heap-in-go-6e3

package dijkstra

type Path struct {
	value int      `json:"-"`
	nodes []string `json:"-"`
}

type minpath []Path
type Heap struct {
	values *minpath `json:"-"`
}

type Edge struct {
	node   string `json:"-"`
	weight int    `json:"-"`
}

type Graph struct {
	nodes map[string][]Edge `json:"-"`
}
