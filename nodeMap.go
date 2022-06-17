package UrlCrawler

type nodeMap struct {
	m                map[nodeIndex]int
	name_buffer      *nameMapper
	cachedFirst      nodeName
	cachedFirstValid bool
}

func newNodeMap(nMapper *nameMapper) *nodeMap {
	nm := &nodeMap{make(map[nodeIndex]int), nMapper, "", false}
	return nm
}

func (y *nodeMap) append(x nodeName) {
	y.cachedFirstValid = false
	idx, ok := y.name_buffer.getIdx(x)
	if !ok {
		idx = y.name_buffer.append(x)
	}
	y.m[idx]++
}

func (y *nodeMap) delete(x nodeName) {
	y.cachedFirstValid = false
	if v, ok := y.name_buffer.getIdx(x); ok {
		delete(y.m, v)
	}
}

// func (y *nodeMap) isEmpty() bool {
// 	return len(y.m) > 0
// }

func (y *nodeMap) getFirst() (nodeName, bool) {
	if y.cachedFirstValid {
		return y.cachedFirst, true
	}
	for k := range y.m {
		if k, e := y.name_buffer.getName(k); e {
			y.cachedFirst = k
			y.cachedFirstValid = true
			return k, true
		}
	}
	return "", false
}

func (y *nodeMap) isExist(name nodeName) bool {
	if idx, ok := y.name_buffer.getIdx(name); ok {
		if _, b := y.m[idx]; b {
			return true
		}
	}
	return false
}
