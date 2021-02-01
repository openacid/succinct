package main

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/google/btree"
	"github.com/openacid/low/mathext/zipf"
	"github.com/openacid/low/size"
	"github.com/openacid/succinct"
	"github.com/openacid/testkeys"
)

var OutputSet bool

func main() {
	for _, name := range []string{
		"200kweb2",
		"870k_ip4_hex",
	} {
		keys := getKeys(name)
		oriSize := 0
		for _, k := range keys {
			oriSize += len(k)
		}
		size := 0
		nsop := 0
		var rst testing.BenchmarkResult

		outp := func(typ string, size, nsop int) {
			fmt.Printf("%s %s: size: %dkb %d%% ns/op: %d\n", name,
				typ, size/1024, size*100/oriSize, nsop)
		}

		{
			rst = testing.Benchmark(func(b *testing.B) {
				size = subBenArrayBsearch(b, keys)
			})
			nsop = int(rst.NsPerOp())
			outp("bsearch", size, nsop)
		}

		{
			rst = testing.Benchmark(func(b *testing.B) {
				size = subBenSetHas(b, keys)
			})
			nsop = int(rst.NsPerOp())
			outp("succinct.Set", size, nsop)
		}

		{
			rst = testing.Benchmark(func(b *testing.B) {
				size = subBenBtree(b, keys)
			})
			nsop = int(rst.NsPerOp())
			outp("btree", size, nsop)
		}
	}

}

func subBenSetHas(b *testing.B, keys []string) int {
	sz := 0
	b.Run(fmt.Sprintf("Has:n=%d", len(keys)), func(b *testing.B) {
		s := succinct.NewSet(keys)

		load := zipf.Accesses(2, 1.5, len(keys), b.N, nil)
		sz = size.Of(s)

		b.ResetTimer()

		var v bool
		for i := 0; i < b.N; i++ {
			rst := s.Has(keys[load[i]])
			v = v || rst
		}
		OutputSet = v
	})

	return sz
}

func subBenArrayBsearch(b *testing.B, keys []string) int {
	sz := 0
	b.Run(fmt.Sprintf("bsearch:n=%d", len(keys)), func(b *testing.B) {

		load := zipf.Accesses(2, 1.5, len(keys), b.N, nil)
		sz = size.Of(keys)

		b.ResetTimer()

		var v bool
		for i := 0; i < b.N; i++ {
			rst := sort.SearchStrings(keys, keys[load[i]])
			v = v || rst == 8
		}
		OutputSet = v
	})

	return sz
}

type BtreeElt struct {
	Key string
}

func (kv *BtreeElt) Less(than btree.Item) bool {
	o := than.(*BtreeElt)
	return kv.Key < o.Key
}

func subBenBtree(b *testing.B, keys []string) int {
	sz := 0
	b.Run(fmt.Sprintf("btree:n=%d", len(keys)), func(b *testing.B) {

		bt := btree.New(32)
		for _, k := range keys {
			v := &BtreeElt{Key: k}
			bt.ReplaceOrInsert(v)
		}

		load := zipf.Accesses(2, 1.5, len(keys), b.N, nil)
		sz = size.Of(bt)

		b.ResetTimer()

		var id int
		for i := 0; i < b.N; i++ {
			idx := load[i]
			itm := &BtreeElt{Key: keys[idx]}
			ee := bt.Get(itm)
			id += len(ee.(*BtreeElt).Key)
		}
		OutputSet = id > 8
	})
	return sz
}

func getKeys(name string) []string {

	keys := testkeys.Load(name)
	if name == "870k_ip4_hex" {
		ks := make([]string, 0, len(keys))
		for _, k := range keys {
			n, err := strconv.ParseUint(k, 16, 0)
			if err != nil {
				panic(err)
			}

			packed := string([]byte{
				byte(n >> 24),
				byte(n >> 16),
				byte(n >> 8),
				byte(n),
			})

			ks = append(ks, packed)
		}
		return ks
	}
	return keys
}
