package UrlCrawler

import "testing"

func TestNodeMap(t *testing.T) {

	nm := newNodeMap(newNameMapper(128))

	nm.delete("A")

	_, ok := nm.getFirst()
	if ok != false {
		t.Log("Result: ", ok)
		t.Error("Must be:", false)
	}
	ok = nm.isExist("A")
	if ok != false {
		t.Log("Result: ", ok)
		t.Error("Must be:", false)
	}

	nm.append("A")
	nm.append("B")
	ok = nm.isExist("A")
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}
	ok = nm.isExist("B")
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}
	ok = nm.isExist("C")
	if ok != false {
		t.Log("Result: ", ok)
		t.Error("Must be:", false)
	}
	_, ok = nm.getFirst()
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}

	nm.delete("A")
	ok = nm.isExist("A")
	if ok != false {
		t.Log("Result: ", ok)
		t.Error("Must be:", false)
	}
	ok = nm.isExist("B")
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}
	v, ok := nm.getFirst()
	if v == "A" {
		t.Log("Result: ", v)
		t.Error("Must be:", "B")
	}
	if v != "B" {
		t.Log("Result: ", v)
		t.Error("Must be:", "B")
	}
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}

	nm = newNodeMap(newNameMapper(128))
	nm.append("A")
	nm.append("B")
	nm.delete("B")
	ok = nm.isExist("B")
	if ok != false {
		t.Log("Result: ", ok)
		t.Error("Must be:", false)
	}
	ok = nm.isExist("A")
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}
	_, ok = nm.getFirst()
	if ok != true {
		t.Log("Result: ", ok)
		t.Error("Must be:", true)
	}

}
