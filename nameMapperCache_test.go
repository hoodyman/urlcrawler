package UrlCrawler

import (
	"math"
	"testing"
)

func TestNameMapperCache(t *testing.T) {

	var cache_size = 128

	cache := newNameMapperCache(cache_size)

	success := cache.getCacheSuccess()
	if success != 0 {
		t.Log("Result:", success)
		t.Error("Must be 0")
	}

	node_name, node_idx := cache.getName(0)
	if node_name != "" {
		t.Log("Result:", node_name)
		t.Error("Must be \"\"")
	}
	if node_idx != false {
		t.Log("Result:", node_idx)
		t.Error("Must be False")
	}

	success = cache.getCacheSuccess()
	if success != 0 {
		t.Log("Result:", success)
		t.Error("Must be 0")
	}

	cache.reset()
	cache.putName("A", 0)

	success = cache.getCacheSuccess()
	if success != 0 {
		t.Log("Result:", success)
		t.Error("Must be 0")
	}

	node_name, node_idx = cache.getName(0)
	if node_name != "A" {
		t.Log("Result:", node_name)
		t.Error("Must be 1")
	}
	if node_idx != true {
		t.Log("Result:", node_idx)
		t.Error("Must be True")
	}
	success = cache.getCacheSuccess()
	if success != 1 {
		t.Log("Result:", success)
		t.Error("Must be 1")
	}

	node_name, node_idx = cache.getName(1)
	if node_name != "" {
		t.Log("Result:", node_name)
		t.Error("Must be \"\"")
	}
	if node_idx != false {
		t.Log("Result:", node_idx)
		t.Error("Must be False")
	}
	success = cache.getCacheSuccess()
	if success != 0.5 {
		t.Log("Result:", success)
		t.Error("Must be 0.5")
	}

	cache = newNameMapperCache(2)

	cache.putName("A", 0)
	cache.putName("C", 2)
	cache.putName("B", 1)
	cache.putName("B", 1) // by idx repeat

	node_name, node_idx = cache.getName(1)
	if node_name != "B" {
		t.Log("Result:", node_name)
		t.Error("Must be B")
	}
	if node_idx != true {
		t.Log("Result:", node_idx)
		t.Error("Must be True")
	}
	success = cache.getCacheSuccess()
	if success != 1 {
		t.Log("Result:", success)
		t.Error("Must be 1")
	}

	// test success ocerflow
	target_relative_error := 1e-6
	// 1_000 hits
	cache.reset()
	cache.putName("A", 0)
	cache.getName(1) // cache miss
	n := 1_000
	for i := 0; i < n; i++ {
		cache.getName(0)
	}
	success = cache.getCacheSuccess()
	true_success := float64(n) / float64(n+1)
	relative_error := math.Abs(success-true_success) / true_success
	if relative_error > target_relative_error {
		t.Logf("Result relative error: %e", relative_error)
		t.Errorf("Must be relative error: %e", target_relative_error)
	}
	// 1_000_000 hits
	cache.reset()
	cache.putName("A", 0)
	cache.getName(1) // cache miss
	n = 1_000_000
	for i := 0; i < n; i++ {
		cache.getName(0)
	}
	success = cache.getCacheSuccess()
	true_success = float64(n) / float64(n+1)
	relative_error = math.Abs(success-true_success) / true_success
	if relative_error > target_relative_error {
		t.Logf("Result relative error: %e", relative_error)
		t.Errorf("Must be relative error: %e", target_relative_error)
	}

}
