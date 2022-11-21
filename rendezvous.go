// Package rendezvous provides a simple concurrent implementation of Weighted Rendezvous Hashing.
// Rendezvous Hashing is also known as Highest Random Weight (HRW) hashing.
//
// [Disclaimer]
// This package is not meant to be used in production!
// I made it for learning purposes. It hasn't been thoroughly tested and benchmarked.
//
// Here're some questions to answer & some tasks to do before it can be production ready:-
// 1. Is there a better, faster way to compute weighted scores?
// 2. The package uses FNV1-a by default, could Murmur be a better hashing function?
// 3. How does it perform? It needs to be benchmarked.
// 4. Is sync.Map a better choice than map + sync.RWMutex?
//
// [Example]
//
//      r := rendezvous.New()
//      r.AddWeightedNodes(map[string]int{
//          "s1.test.com": 5,
//          "s2.test.com": 1,
//          "s3.test.com": 10,
//      })
//
//      selected := r.Resolve("some-request")
//      selected = r.Resolve("some-other-request")
//
// References:-
// https://en.wikipedia.org/wiki/Rendezvous_hashing
// https://www.snia.org/sites/default/files/SDC15_presentations/dist_sys/Jason_Resch_New_Consistent_Hashings_Rev.pdf
//

package rendezvous

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"sync"
)



type rendezvous struct {
	Hasher hash.Hash32 // 32-bit hash generator that implements hash.Hash32 interface
	// Bring your own 32-bit hasher implementation if you want to

	nodes map[string]float64 // keeps track of node-weight pairs
	mu    sync.RWMutex       // protects nodes
}

// New creates a new instances of rendezvous
func New() *rendezvous {
	return &rendezvous{
		Hasher: fnv.New32a(),
		nodes:  map[string]float64{},
	}
}

// AddNode adds a node with weight=1
func (r *rendezvous) AddNode(addr string) {
	r.AddWeightedNode(addr, 1)
}

func (r *rendezvous) AddWeightedNode(addr string, weight int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[addr] = float64(weight)
}

func (r *rendezvous) AddWeightedNodes(pairs map[string]int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, v := range pairs {
		r.nodes[k] = float64(v)
	}
}

func (r *rendezvous) RemoveNode(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, _ := range r.nodes {
		if k == addr {
			delete(r.nodes, k)
			return
		}
	}
}

// Len returns the number of nodes
func (r *rendezvous) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.nodes)
}

// Resolve returns the node that has the highest score
func (r *rendezvous) Resolve(request string) string {
    if len(r.nodes) == 0 {
        // Good time to panic because it indicates logic error on the developer's part
        // Convert this code to return an error if it's normal to have zero active servers
        panic("No server available")
    }

	node := ""
	var maxScore float64 = 0.0

	for addr := range r.nodes {
		s := r.computeWeightedScore(addr, request)
		if s > maxScore {
			maxScore = s
			node = addr
		}
	}

	return node
}

// computerWeightedScore calculates the weighted score for a node-request pair
func (r *rendezvous) computeWeightedScore(addr string, request string) float64 {
	r.mu.RLock()
	weight, found := r.nodes[addr]
	r.mu.RUnlock()

	if !found {
		panic(fmt.Sprintf("%v is not a valid node.", addr))
	}

	score := r.normalize(r.hash(fmt.Sprintf("%v:%v", addr, request)))

	// See the references in the package description for more details about this formula
	logScore := 1.0 / -math.Log(score)

	return weight * logScore
}

// Equivalent to float64(math.MaxUint32) + 1.0
const uint32Bound = float64(1 << 32)

// normalize normalizes a uint32 number to the (0,1] range
func (r *rendezvous) normalize(a uint32) float64 {
	// Convert a from uint32 to float64 for two reasons:
	// 1. return type is float64
	// 2. so that the result doesn't overflow when we add 1 to it.
	//
	// Add 1 to float64(a) to make sure that its never 0

	return (float64(a) + 1.0) / uint32Bound
}

// hash generates a 32-bit hash based on the input
func (r *rendezvous) hash(key string) uint32 {
	r.Hasher.Reset()
	r.Hasher.Write([]byte(key))
	return r.Hasher.Sum32()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
