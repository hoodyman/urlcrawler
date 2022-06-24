package UrlCrawler

type nameMapper struct {
	m     map[nodeName]nodeIndex
	cache *nameMapperCache
}

func newNameMapper(cache_size int) *nameMapper {
	name_mapper := &nameMapper{}
	name_mapper.m = make(map[nodeName]nodeIndex)
	name_mapper.cache = newNameMapperCache(cache_size)
	return name_mapper
}

func (y *nameMapper) len() int {
	return len(y.m)
}

func (y *nameMapper) append(x nodeName) nodeIndex {
	if v, ok := y.m[x]; ok {
		return v
	} else {
		idx := len(y.m)
		y.m[x] = nodeIndex(idx)
		return nodeIndex(idx)
	}
}

func (y *nameMapper) getIdx(x nodeName) (nodeIndex, bool) {
	if v, ok := y.m[x]; ok {
		return v, true
	} else {
		return 0, false
	}
}

func (y *nameMapper) getName(x nodeIndex) (nodeName, bool) {
	if v, ok := y.cache.getName(x); ok {
		return v, true
	}
	if x < nodeIndex(len(y.m)) {
		for k, v := range y.m {
			if v == x {
				y.cache.putName(k, x)
				return k, true
			}
		}
	}
	return "", false
}

func (y *nameMapper) getIdxOrAppend(x nodeName) nodeIndex {
	if a, b := y.getIdx(x); b {
		return a
	}
	return y.append(x)
}

func (y *nameMapper) getCacheSuccess() float64 {
	return y.cache.getCacheSuccess()
}
