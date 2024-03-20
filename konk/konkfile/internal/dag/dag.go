package dag

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// A DAG is a directed acyclic graph.
type DAG[Node comparable] struct {
	nodes map[Node][]Node
}

// Nodes returns the nodes in the DAG.
func (d *DAG[Node]) Nodes() []Node {
	return maps.Keys(d.nodes)
}

// From returns the nodes that the given node points to.
func (d *DAG[Node]) From(n Node) []Node {
	return d.nodes[n]
}

// AddNode adds a new node to the DAG.
func (d *DAG[Node]) AddNode(n Node) {
	if d.hasNode(n) {
		return
	}

	d.nodes[n] = []Node{}
}

// AddEdge adds a new edge to the DAG.
func (d *DAG[Node]) AddEdge(from Node, to Node) error {
	if !d.hasNode(from) {
		return errors.New("from node does not exist")
	}

	if !d.hasNode(to) {
		return errors.New("to node does not exist")
	}

	if d.hasEdge(from, to) {
		return nil
	}

	d.nodes[from] = append(d.nodes[from], to)

	return nil
}

func (d *DAG[Node]) hasNode(n Node) bool {
	_, ok := d.nodes[n]
	return ok
}

func (d *DAG[Node]) hasEdge(from Node, to Node) bool {
	if !d.hasNode(from) {
		return false
	}

	return slices.Contains(d.nodes[from], to)
}

func (d *DAG[Node]) Visit(n Node) ([]Node, error) {
	if !d.hasNode(n) {
		return nil, errors.New("node does not exist")
	}

	visited := make([]Node, 0, len(d.nodes))
	result := make([]Node, 0, len(d.nodes))

	result, err := d.visit(n, visited, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DAG[Node]) visit(n Node, visited []Node, result []Node) ([]Node, error) {
	if slices.Contains(visited, n) {
		return nil, fmt.Errorf("cycle detected: %v in %v", n, visited)
	}

	visited = append(visited, n)

	for _, m := range d.nodes[n] {
		subresult, err := d.visit(m, slices.Clone(visited), result)
		if err != nil {
			return nil, err
		}

		result = subresult
	}

	if !slices.Contains(result, n) {
		result = append(result, n)
	}

	return result, nil
}

// New creates a new DAG.
func New[Node comparable]() *DAG[Node] {
	return &DAG[Node]{
		nodes: make(map[Node][]Node),
	}
}
