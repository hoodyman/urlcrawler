package UrlCrawler

import (
	"math"
	"testing"
)

func TestNameMapper(t *testing.T) {

	mapper := newNameMapper(2)

	_, ok := mapper.getIdx("A")
	if ok != false {
		t.Log("Result:", ok)
		t.Error("Must be: False")
	}

	n, ok := mapper.getName(0)
	if ok != false {
		t.Log("Result:", ok)
		t.Error("Must be: False")
	}
	if n != "" {
		t.Log("Result:", n)
		t.Error("Must be: \"\"")
	}

	idx := mapper.append("A")
	if idx != 0 {
		t.Log("Result:", idx)
		t.Error("Must be: 0")
	}
	idx = mapper.append("A") // value repeat
	if idx != 0 {
		t.Log("Result:", idx)
		t.Error("Must be: 0")
	}

	idx, ok = mapper.getIdx("A")
	if ok != true {
		t.Log("Result:", ok)
		t.Error("Must be: True")
	}
	if idx != 0 {
		t.Log("Result:", idx)
		t.Error("Must be: 0")
	}

	n, ok = mapper.getName(0)
	if n != "A" {
		t.Log("Result:", n)
		t.Error("Must be: A")
	}
	if ok != true {
		t.Log("Result:", ok)
		t.Error("Must be: True")
	}

	mapper = newNameMapper(2)
	mapper.append("A")
	n, ok = mapper.getName(0)
	if n != "A" {
		t.Log("Result:", n)
		t.Error("Must be: A")
	}
	if ok != true {
		t.Log("Result:", ok)
		t.Error("Must be: True")
	}
	n, ok = mapper.getName(0)
	if n != "A" {
		t.Log("Result:", n)
		t.Error("Must be: A")
	}
	if ok != true {
		t.Log("Result:", ok)
		t.Error("Must be: True")
	}
	n, ok = mapper.getName(0)
	if n != "A" {
		t.Log("Result:", n)
		t.Error("Must be: A")
	}
	if ok != true {
		t.Log("Result:", ok)
		t.Error("Must be: True")
	}

	success := mapper.getCacheSuccess()
	true_success := float64(2) / float64(3)
	relative_error := math.Abs(success-true_success) / true_success
	if relative_error > 1e-6 {
		t.Log(success)
		t.Logf("Result relative error: %e", relative_error)
		t.Errorf("Must be relative error: %e", 1e-6)
	}
}
