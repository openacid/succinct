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
	//   Compressed size: 1209 KB, ratio: 54 %
	// Memory layout:
	// *succinct.Set: 1238864
	//     succinct.Set: 1238856
	//         leaves: []uint64: 99128
	//             0: uint64: 8
	//         labelBitmap: []uint64: 198224
	//             0: uint64: 8
	//         labels: []uint8: 792800
	//             0: uint8: 1
	//         ranks: []int32: 99128
	//             0: int32: 4
	//         selects: []int32: 49576
	//             0: int32: 4
}
