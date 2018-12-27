package main

import "fmt"

type Node struct {
	id      int
	visited bool
}

func (n *Node) String() string {
	return fmt.Sprintf("%v %t", n.id, n.visited)
}

type Graph struct {
	nodes []*Node
	edges map[Node][]*Node
}

func (g *Graph) AddNode(n *Node) {
	g.nodes = append(g.nodes, n)
}

func (g *Graph) AddEdge(n1, n2 *Node) {
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
}

func (g *Graph) String() string {
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	return s
}

func (g *Graph) GetEdges(n *Node) []int {
	e := g.edges[*n]
	var paths []int
	for _, path := range e {
		paths = append(paths, path.id)
	}
	return paths
}
