package UrlCrawler

import (
	"math"
)

// Call Init()
type nameMapperCache struct {
	cacheHit             float64
	cacheMiss            float64
	statsCoefficient     float64
	statsCoefficientMult float64
	supressCasheUpdate   bool
	cacheSize            int
	cache                []cacheLine
	currentAmount        int
}

type cacheLine struct {
	name  nodeName
	index nodeIndex
}

const defaultStatsCoefficient = 1
const defaultStatsCoefficientMult = 0.001

func newNameMapperCache(cache_size int) *nameMapperCache {
	nbc := &nameMapperCache{}
	nbc.cacheSize = cache_size
	nbc.reset()
	return nbc
}

func (y *nameMapperCache) reset() {
	y.cacheHit = 0
	y.cacheMiss = 0
	y.statsCoefficient = defaultStatsCoefficient
	y.statsCoefficientMult = defaultStatsCoefficientMult
	y.supressCasheUpdate = false
	y.cache = make([]cacheLine, y.cacheSize)
	y.currentAmount = 0
}

func (y *nameMapperCache) getName(idx nodeIndex) (nodeName, bool) {
	for i, v := range y.cache {
		if i < y.currentAmount && v.index == idx {
			if !y.supressCasheUpdate {
				if i > 0 {
					y.cache[i], y.cache[i-1] = y.cache[i-1], y.cache[i]
				}
				y.cacheHit += y.statsCoefficient
				y.statsOverflowCheck()
			}
			return v.name, true
		}
	}
	if !y.supressCasheUpdate {
		y.cacheMiss += y.statsCoefficient
		y.statsOverflowCheck()
	}
	return "", false
}

// func (y *nameMapperCache) getIndex(name nodeName) (nodeIndex, bool) {
// 	for i, v := range y.cache {
// 		if i < y.currentAmount && v.name == name {
// 			if !y.supressCasheUpdate {
// 				if i > 0 {
// 					y.cache[i], y.cache[i-1] = y.cache[i-1], y.cache[i]
// 				}
// 				y.cacheHit += y.statsCoefficient
// 				y.statsOverflowCheck()
// 			}
// 			return v.index, true
// 		}
// 	}
// 	if !y.supressCasheUpdate {
// 		y.cacheMiss += y.statsCoefficient
// 		y.statsOverflowCheck()
// 	}
// 	return 0, false
// }

func (y *nameMapperCache) putName(name nodeName, idx nodeIndex) {
	y.supressCasheUpdate = true
	defer func() { y.supressCasheUpdate = false }()
	if _, ok := y.getName(idx); ok {
		return
	}
	if y.currentAmount == y.cacheSize {
		y.cache[y.currentAmount-1] = cacheLine{name, idx}
	} else {
		y.cache[y.currentAmount] = cacheLine{name, idx}
		y.currentAmount++
	}
}

func (y *nameMapperCache) getCacheSuccess() float64 {
	x := y.cacheHit + y.cacheMiss
	if x == 0 {
		return 0
	}
	return y.cacheHit / x
}

func (y *nameMapperCache) statsOverflowCheck() {
	max := math.Max(y.cacheHit, y.cacheMiss)
	if max >= 1_000 { //math.MaxFloat64-1 {
		y.statsCoefficient *= y.statsCoefficientMult
		y.cacheHit *= y.statsCoefficient
		y.cacheMiss *= y.statsCoefficient
	}
}
