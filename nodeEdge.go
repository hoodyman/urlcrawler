package UrlCrawler

type nodeEdge struct {
	NodeA nodeIndex
	NodeB nodeIndex
}

func newNodeEdge(a, b nodeName, name_buffer *nameMapper) nodeEdge {
	x, y := name_buffer.append(a), name_buffer.append(b)
	if x < y {
		return nodeEdge{x, y}
	}
	return nodeEdge{y, x}
}

type nodeEdgeMap map[nodeEdge]weightType

// return true if value is in container
func (y *nodeEdgeMap) append(x nodeEdge) bool {
	_, n := (*y)[x]
	(*y)[x]++
	return n
}

func (y *nodeEdgeMap) len() int {
	return len(*y)
}
