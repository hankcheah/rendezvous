package rendezvous

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	r := New()
	r.AddNode("w1 node")

	// This will overwrite the first node
	r.AddWeightedNode("w1 node", 2)
	r.AddWeightedNode("w50 node", 50)

	if r.Len() != 2 {
		t.Errorf("Expected 2 nodes.")
	}

	expected := map[string]float64{"w1 node": 2, "w50 node": 50}
	if !reflect.DeepEqual(r.nodes, expected) {
		t.Errorf("Expected: %v Got %v", expected, r.nodes)
	}
}

func TestAddWeightedNodes(t *testing.T) {
	nodes := map[string]int{
		"s1.test.net": 5,
		"s2,test.net": 4,
		"s3.test.net": 3,
		"s4.test.net": 2,
	}

	r := New()
	r.AddWeightedNodes(nodes)

	if r.Len() != 4 {
		t.Errorf("Expected 2 nodes.")
	}

	expected := map[string]float64{
		"s1.test.net": 5.0,
		"s2,test.net": 4.0,
		"s3.test.net": 3.0,
		"s4.test.net": 2.0,
	}
	if !reflect.DeepEqual(r.nodes, expected) {
		t.Errorf("Expected: %v Got %v", expected, r.nodes)
	}
}

func TestRemove(t *testing.T) {
	r := New()

	r.AddNode("w1 node")
	r.AddWeightedNode("w50 node", 50)

	if r.Len() != 2 {
		t.Errorf("Expected 2 nodes.")
	}

	r.RemoveNode("w1 node")

	expected := map[string]float64{"w50 node": 50}
	if !reflect.DeepEqual(r.nodes, expected) {
		t.Errorf("Expected: %v Got %v", expected, r.nodes)
	}
}

// Just a straightfoward test. Test result might change if the hash function changes.
func TestResolve(t *testing.T) {
	r := New()

	r.AddNode("s1.test.com")
	r.AddNode("s2.test.com")

	selected := r.Resolve("request")
	expected := "s2.test.com"
	if selected != expected {
		t.Errorf("Expected: %v Got %v", expected, selected)
	}
}

// Node with higher score gets picked over other nodes with lower scores
func TestResolveWeighted(t *testing.T) {
	r := New()

	r.AddWeightedNode("s1.test.com", 10)
	r.AddWeightedNode("s2.test.com", 2)
	r.AddWeightedNode("s3.test.com", 3)

	selected := r.Resolve("request")
	expected := "s2.test.com"
	if selected != expected {
		t.Logf("%v\n", r.nodes)
		t.Errorf("Expected: %v Got %v", expected, selected)
	}
}
