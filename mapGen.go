package main

import (
	"math/rand"
	"time"
)

type MapGen struct {
	Graph     Graph
	generated bool
}

func NewMapGen() MapGen {
	return MapGen{
		Graph: Graph{},
	}
}

func (mg *MapGen) GenerateMap(nodeCount, moreEdgesChance, edgesMin, edgesMax int) {
	rand.Seed(time.Now().UnixNano())
	mg.generated = true
	for i := 0; i < nodeCount; i++ {
		mg.Graph.AddNode(&Node{
			id:      i,
			visited: false,
		})
	}

	for i := 0; i < nodeCount; i++ {
		to := mg.genEdge(i, nodeCount)
		mg.Graph.AddEdge(
			mg.Graph.nodes[i],
			mg.Graph.nodes[to])
		if Chancer(moreEdgesChance) == true {
			moreEdgesCount := Randomize(edgesMin, edgesMax)
			for j := 0; j < moreEdgesCount; j++ {
				to := mg.genEdge(i, nodeCount)
				mg.Graph.AddEdge(
					mg.Graph.nodes[i],
					mg.Graph.nodes[to])
			}
		}
	}

	unreachables := mg.CheckReachability()

	for _, unreachable := range unreachables {
		for {
			to := mg.genEdge(unreachable.id, nodeCount)
			if mg.Graph.nodes[to].visited == true {
				mg.Graph.AddEdge(
					mg.Graph.nodes[unreachable.id],
					mg.Graph.nodes[to],
				)
				unreachable.visited = true
				break
			}
		}
	}

	for _, nod := range mg.Graph.nodes {
		nod.visited = false
	}

}

func (mg *MapGen) genEdge(from, count int) int {
	for {
		rnd := Randomize(0, count)
		if rnd != from && mg.CheckMap(from) == false {
			return rnd
		}
	}
}

func (mg *MapGen) CheckMap(node int) bool {
	edges, exists := mg.Graph.edges[*mg.Graph.nodes[node]]
	if exists == true {
		for i := 0; i < len(edges); i++ {
			if edges[i].id == mg.Graph.nodes[node].id {
				return true
			}
		}
	}
	return false
}

func (mg *MapGen) CheckReachability() []*Node {
	visits := NewStack()
	edgesFirstNode := mg.Graph.edges[*mg.Graph.nodes[0]]
	mg.Graph.nodes[0].visited = true
	for _, edge := range edgesFirstNode {
		visits.Push(edge.id)
	}

	for visits.Len() > 0 {
		visiting := visits.Pop()
		if mg.Graph.nodes[visiting].visited == false {
			for _, edge := range mg.Graph.edges[*mg.Graph.nodes[visiting]] {
				visits.Push(edge.id)
			}
		}
		mg.Graph.nodes[visiting].visited = true
	}

	unvisited := make([]*Node, 0)

	for _, node := range mg.Graph.nodes {
		if node.visited == false {
			unvisited = append(unvisited, node)
		}
	}

	return unvisited
}

func Chancer(chance int) bool {
	rnd := Randomize(0, 100)
	if chance == 0 {
		return false
	}
	if rnd < chance {
		return true
	}
	return false
}

func Randomize(min, max int) int {
	rand.Intn(max)
	rest := max - min
	rnd := rand.Intn(rest)
	return rnd + min
}
