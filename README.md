# succinct

[![Travis](https://travis-ci.com/openacid/succinct.svg?branch=main)](https://travis-ci.com/openacid/succinct)
![test](https://github.com/openacid/succinct/workflows/test/badge.svg)

[![Report card](https://goreportcard.com/badge/github.com/openacid/succinct)](https://goreportcard.com/report/github.com/openacid/succinct)
[![Coverage Status](https://coveralls.io/repos/github/openacid/succinct/badge.svg?branch=main&service=github)](https://coveralls.io/github/openacid/succinct?branch=main&service=github)

[![GoDoc](https://godoc.org/github.com/openacid/succinct?status.svg)](http://godoc.org/github.com/openacid/succinct)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/openacid/succinct)](https://pkg.go.dev/github.com/openacid/succinct)
[![Sourcegraph](https://sourcegraph.com/github.com/openacid/succinct/-/badge.svg)](https://sourcegraph.com/github.com/openacid/succinct?badge)

succinct provides several static succinct data types

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Succinct Set](#succinct-set)
  - [Implementation](#implementation)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Succinct Set

中文介绍: [100行代码的压缩前缀树: 50% smaller](https://blog.openacid.com/algo/succinctset/)

Set is a succinct, sorted and static string set impl with compacted trie as
storage. The space cost is about half lower than the original data.

```go
package succinct

import "fmt"

func ExampleNewSet() {
	keys := []string{
		"A", "Aani", "Aaron", "Aaronic", "Aaronical", "Aaronite",
		"Aaronitic", "Aaru", "Ab", "Ababdeh", "Ababua", "Abadite",
	}
	s := NewSet(keys)
	for _, k := range []string{"Aani", "Foo", "Ababdeh"} {
		found := s.Has(k)
		fmt.Printf("lookup %10s, found: %v\n", k, found)
	}

	// Output:
	//
	// lookup       Aani, found: true
	// lookup        Foo, found: false
	// lookup    Ababdeh, found: true
}
```

A benchmark with 200 kilo real-world words collected from web shows that:
- the space a `Set` costs is only **57%** of original data size.
- And a `Has()` costs about `400 ns` with a **zip-f** workload.

```go
package succinct

import (
	"fmt"

	"github.com/openacid/low/size"
	"github.com/openacid/testkeys"
)

func ExampleNewSet_memory() {

	keys := testkeys.Load("200kweb2")

	original := 0
	for _, k := range keys {
		original += len(k)
	}

	s := NewSet(keys)

	fmt.Println("With", len(keys)/1000, "thousands keys:")
	fmt.Println("  Original size:", original/1024, "KB")
	fmt.Println("  Compressed size:",
		size.Of(s)/1024, "KB, ratio:",
		size.Of(s)*100/original, "%")
	fmt.Println("Memory layout:")
	fmt.Println(size.Stat(s, 10, 1))

	// Output:
	//
	// With 235 thousands keys:
	//   Original size: 2204 KB
	//   Compressed size: 1258 KB, ratio: 57 %
	// Memory layout:
	// *succinct.Set: 1288412
	//     succinct.Set: 1288404
	//         leaves: []uint64: 99128
	//             0: uint64: 8
	//         labelBitmap: []uint64: 198224
	//             0: uint64: 8
	//         labels: []uint8: 792800
	//             0: uint8: 1
	//         ranks: []int32: 99128
	//             0: int32: 4
	//         selects: []int32: 99124
	//             0: int32: 4
}
```

## Implementation

It stores sorted strings in a compacted trie(AKA prefix tree). A trie node has
at most 256 outgoing labels. A label is just a single byte. E.g., [ab, abc,
abcd, axy, buv] is represented with a trie like the following: (Numbers are node
id)

    ^ -a-> 1 -b-> 3 $
      |      |      `c-> 6 $
      |      |             `d-> 9 $
      |      `x-> 4 -y-> 7 $
      `b-> 2 -u-> 5 -v-> 8 $

Internally it uses a packed []byte and a bitmap with `len([]byte)` bits to
describe the outgoing labels of a node,:

    ^: ab  00
    1: bx  00
    2: u   0
    3: c   0
    4: y   0
    5: v   0
    6: d   0
    7: ø
    8: ø
    9: ø

In storage it packs labels together and bitmaps joined with separator `1`:

    labels(ignore space): "ab bx u c y v d"
    label bitmap:          0010010101010101111

Finally leaf nodes are indicated by another bitmap `leaves`, in which a `1` at
i-th bit indicates the i-th node is a leaf:

    leaves: 0001001111

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.