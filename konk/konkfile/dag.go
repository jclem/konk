package konkfile

import (
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type edge[node comparable] struct {
	from node
	to   node
}

type dag[node comparable] struct {
	nodes []node
	edges []edge[node]
}

func (d *dag[node]) addNode(n node) {
	if d.hasNode(n) {
		return
	}

	d.nodes = append(d.nodes, n)
}

func (d *dag[node]) addEdge(e edge[node]) error {
	if !d.hasNode(e.from) {
		return errors.New("from node does not exist")
	}

	if !d.hasNode(e.to) {
		return errors.New("to node does not exist")
	}

	for _, de := range d.edges {
		if de.from == e.from && de.to == e.to {
			return nil
		}
	}

	d.edges = append(d.edges, e)

	return nil
}

func (d *dag[node]) hasNode(n node) bool {
	return slices.Contains(d.nodes, n)
}

func (d *dag[node]) from(n node) []node {
	ns := make(map[node]struct{})

	for _, de := range d.edges {
		if de.from == n {
			ns[de.to] = struct{}{}
		}
	}

	return maps.Keys(ns)
}

func (d *dag[node]) to(n node) []node {
	ns := make(map[node]struct{})

	for _, de := range d.edges {
		if de.to == n {
			ns[de.from] = struct{}{}
		}
	}

	return maps.Keys(ns)
}

func newDAG[node comparable]() *dag[node] {
	return &dag[node]{
		nodes: []node{},
		edges: []edge[node]{},
	}
}
