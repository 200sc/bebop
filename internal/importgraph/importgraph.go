package importgraph

import "fmt"

type dgraph struct {
	nodes map[string][]string
}

func NewDgraph() *dgraph {
	return &dgraph{
		nodes: make(map[string][]string),
	}
}

func (d *dgraph) AddEdge(from, to string) {
	d.nodes[from] = append(d.nodes[from], to)
}

type nodePath struct {
	this string
	prev *nodePath
}

func (p nodePath) string(prefix string) string {
	s := prefix + "\t" + p.this
	if p.prev == nil {
		return s
	}
	return p.prev.string(s + ", imported by:\n")
}

func (d *dgraph) findCycle(from string, stack map[string]struct{}, visited map[string]struct{}, p *nodePath) *nodePath {
	visited[from] = struct{}{}
	stack[from] = struct{}{}
	for _, to := range d.nodes[from] {
		p2 := &nodePath{
			prev: p,
			this: to,
		}
		if _, ok := stack[to]; ok {
			return p2
		}

		if cycle := d.findCycle(to, stack, visited, p2); cycle != nil {
			return cycle
		}
	}
	delete(stack, from)
	return nil
}

func (d *dgraph) FindCycle() error {
	stack := map[string]struct{}{}
	visited := map[string]struct{}{}
	for node := range d.nodes {
		if _, ok := visited[node]; ok {
			continue
		}
		p := &nodePath{this: node}
		if cycle := d.findCycle(node, stack, visited, p); cycle != nil {
			return fmt.Errorf(cycle.string("import cycle found:\n"))
		}
	}
	return nil
}
