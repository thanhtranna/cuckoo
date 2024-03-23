# Cuckoo Filter


Cuckoo filter is a Bloom filter replacement for approximated set-membership queries. While Bloom filters are well-known space-efficient data structures to serve queries like "if item x is in a set?", they do not support deletion. Their variances to enable deletion (like counting Bloom filters) usually require much more space.

Cuckoo ﬁlters provide the ﬂexibility to add and remove items dynamically. A cuckoo filter is based on cuckoo hashing (and therefore named as cuckoo filter). It is essentially a cuckoo hash table storing each key's fingerprint. Cuckoo hash tables can be highly compact, thus a cuckoo filter could use less space than conventional Bloom ﬁlters, for applications that require low false positive rates (< 3%).

For details about the algorithm and citations please use this article for now

["Cuckoo Filter: Better Than Bloom" by Bin Fan, Dave Andersen and Michael Kaminsky](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf)


## Implementation details

The paper cited above leaves several parameters to choose. In this implementation

1. Every element has 2 possible bucket indices
2. Buckets have a static size of 4 fingerprints
3. Fingerprints have a static size of 8 bits

1 and 2 are suggested to be the optimum by the authors. The choice of 3 comes down to the desired false positive rate. Given a target false positive rate of ```r``` and a bucket size ```b```, they suggest choosing the fingerprint size ```f``` using

```
f >= log2(2b/r) bits
```

With the 8 bit fingerprint size in this repository, you can expect ```r ~= 0.03```. [Other implementations](https://github.com/panmari/cuckoofilter) use 16 bit, which correspond to a false positive rate of ```r ~= 0.0001```.

## Example usage:

```go
package main

import (
    "fmt"

    cuckoo "github.com/thanhtranna/cuckoo"
)

func main() {
  cf := cuckoo.NewFilter(1000)
  cf.InsertUnique([]byte("geeky ogre"))

  // Lookup a string (and it a miss) if it exists in the cuckoofilter
  cf.Lookup([]byte("hello"))

  count := cf.Count()
  fmt.Println(count) // count == 1

  // Delete a string (and it a miss)
  cf.Delete([]byte("hello"))

  count = cf.Count()
  fmt.Println(count) // count == 1

  // Delete a string (a hit)
  cf.Delete([]byte("geeky ogre"))

  count = cf.Count()
  fmt.Println(count) // count == 0
  
  cf.Reset()    // reset
}
```

## Benchmarks

```
goos: darwin
goarch: amd64
pkg: github.com/thanhtranna/cuckoo
cpu: Intel(R) Core(TM) i5-9400F CPU @ 2.90GHz
BenchmarkFilter_Reset-6                           626833              1828 ns/op
BenchmarkFilter_Insert-6                         1000000              6519 ns/op
BenchmarkFilter_Lookup-6                         1597498               720.1 ns/op
BenchmarkScalableCuckooFilter_Insert-6           1564302               766.9 ns/op
BenchmarkScalableCuckooFilter_Reset-6             553423              2146 ns/op
BenchmarkScalableCuckooFilter_Lookup-6           1548381               782.8 ns/op
```
