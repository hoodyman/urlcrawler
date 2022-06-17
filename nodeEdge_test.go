package UrlCrawler

import "testing"

func TestNodeEdge(t *testing.T) {

	name_buffer := newNameMapper(128)
	node_edge := newNodeEdge("A", "B", name_buffer)
	t_node_edge := nodeEdge{0, 1}
	if node_edge != t_node_edge {
		t.Log("Result:", node_edge)
		t.Error("Must be:", t_node_edge)
	}

	node_edge = newNodeEdge("B", "A", name_buffer)
	t_node_edge = nodeEdge{0, 1}
	if node_edge != t_node_edge {
		t.Log("Result:", node_edge)
		t.Error("Must be:", t_node_edge)
	}

	name_buffer = newNameMapper(128)
	node_edge = newNodeEdge("B", "A", name_buffer)
	t_node_edge = nodeEdge{0, 1}
	if node_edge != t_node_edge {
		t.Log("Result:", node_edge)
		t.Error("Must be:", t_node_edge)
	}

	name_buffer = newNameMapper(128)
	node_edge = newNodeEdge("A", "B", name_buffer)
	t_node_edge = nodeEdge{0, 1}
	if node_edge != t_node_edge {
		t.Log("Result:", node_edge)
		t.Error("Must be:", t_node_edge)
	}
	node_edge = newNodeEdge("A", "C", name_buffer)
	t_node_edge = nodeEdge{0, 2}
	if node_edge != t_node_edge {
		t.Log("Result:", node_edge)
		t.Error("Must be:", t_node_edge)
	}

	name_buffer = newNameMapper(128)
	node_edge = newNodeEdge("A", "B", name_buffer)
	node_edge_map := nodeEdgeMap{}
	node_edge_map.append(node_edge)
	var n_e nodeEdge
	var res weightType
	for k, v := range node_edge_map {
		n_e = k
		res = v
		break
	}
	if n_e != node_edge {
		t.Log("Result:", n_e)
		t.Error("Must be:", node_edge)
	}
	if res != 1 {
		t.Log("Result:", res)
		t.Error("Must be:", 1)
	}

}
