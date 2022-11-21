# rendezvous
Package rendezvous provides a simple concurrent implementation of Weighted Rendezvous Hashing.
Rendezvous Hashing is originally known as Highest Random Weight (HRW) hashing.

## Examples

### Example 1: Basic Ops
This example shows how to add and remove nodes, as well as how to determine which node is handling a request using resolve().

```
  r := New()

  // Add nodes in bulk
  r.AddWeightedNodes(map[string]int{
    "s1.test.com": 5,
    "s2.test.com": 1,
    "s3.test.com": 10,
  })

  // Add a single node with weight = 1
  r.AddNode("s4.test.com")

  // Add a single weighted node
  r.AddWeightedNode("s5.test.com", 7)
  
  // Remove a node
  r.RemoveNode("s1.test.com")
  
  // Resolve a request to a particular node
  // I.e. select a node
  chosen := r.Resolve("some-request")
```

### Example 2: BYOH - Bring Your Own Hasher
Currently the package only supports 32-bit hash functions and the hasher has to satisfy the hash.Hash32 interface (See https://pkg.go.dev/hash).

```
  import "crc32"

  r.New()
  r.Hasher = crc32.NewIEEE()  // crc32.NewIEEE is a concrete type of hash.Hash32
  
  r.AddNode("node1.test.com")
  selected := r.Resolve("some-request")

```

## Disclaimer

This package is not meant to be used in production!
I made it for learning purposes. It hasn't been thoroughly tested and benchmarked.

Here're some questions to answer & some tasks to do before it can be production ready:-
1. Is there a better, faster way to compute weighted scores?
2. The package uses FNV1-a by default, could Murmur be a better hashing function?
3. How does it perform? It needs to be benchmarked.
4. Is sync.Map a better choice than map + sync.RWMutex?
5. The comments are either too lengthy or not descriptive enough.

## References

Here are some of the references I used to develop this package:-

https://en.wikipedia.org/wiki/Rendezvous_hashing

https://pkg.go.dev/hash

https://github.com/dgryski/go-rendezvous

https://www.snia.org/sites/default/files/SDC15_presentations/dist_sys/Jason_Resch_New_Consistent_Hashings_Rev.pdf
