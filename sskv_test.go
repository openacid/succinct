package succinct

import (
	"sort"
	"testing"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/mathext/zipf"
	"github.com/openacid/testkeys"
	"github.com/openacid/testutil"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	ta := require.New(t)

	type wantTyp struct {
		leaves      string
		labelBitmap string
		ranks       []int32
		selects     []int32
		labels      string
	}

	cases := []struct {
		keys []string
		want wantTyp
	}{
		{
			keys: []string{
				"",
				"a",
			},
			want: wantTyp{
				leaves:      "11000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				labelBitmap: "01100000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				ranks:       []int32{0, 2},
				selects:     []int32{1},
				labels:      "a",
			},
		},
		{
			keys: []string{
				"a",
				"b",
				"c",
			},
			want: wantTyp{
				leaves:      "01110000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				labelBitmap: "00011110 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				ranks:       []int32{0, 4},
				selects:     []int32{3},
				labels:      "abc",
			},
		},
		{
			keys: []string{
				"a",
				"ab",
				"abc",
			},
			want: wantTyp{
				leaves:      "01110000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				labelBitmap: "01010110 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
				ranks:       []int32{0, 4},
				selects:     []int32{1},
				labels:      "abc",
			},
		},
		{
			keys: []string{
				"abc",
				"abcd",
				"abd",
				"abde",
				"bc",
				"bcd",
				"bcde",
				"cde",
			},
			want: wantTyp{
				leaves: "00000101 11111100 00000000 00000000 00000000 00000000 00000000 00000000",
				// 0 -a-> 1 -b-> 4 -c-> 7 -d-> $
				//                 -d-> 8 -e-> $
				//   -b-> 2 -c-> 5 -d-> 9 -e-> $
				//   -c-> 3 -d-> 6 -e-> $
				//
				//        1   2      3       4
				//        abc b c  d cd d e  d e e
				labelBitmap: "00010101 01001010 10101011 11100000 00000000 00000000 00000000 00000000",
				ranks:       []int32{0, 14},
				selects:     []int32{3},
				labels:      "abcbcdcddedee",
			},
		},
		{
			keys: []string{
				"A", "Aani", "Aaron", "Aaronic", "Aaronical", "Aaronite",
				"Aaronitic", "Aaru", "Ab", "Ababdeh", "Ababua", "Abadite",
				"Abama", "Abanic", "Abantes", "Abarambo", "Abaris", "Abasgi",
				"Abassin", "Abatua", "Abba", "Abbadide", "Abbasside", "Abbie",
				"Abby", "Abderian", "Abderite", "Abdiel", "Abdominales", "Abe",
				"Abel", "Abelia", "Abelian", "Abelicea", "Abelite",
				"Abelmoschus", "Abelonian", "Abencerrages", "Aberdeen",
				"Aberdonian", "Aberia", "Abhorson",
			},
			want: wantTyp{
				leaves: ("" +
					"01010000 01010100 00000101 00010001 00010000 00000100 00000000 00101001," +
					"10100010 10000000 10101110 10000010 10000000 10110110 10001011 11001000," +
					"00111010 00000000 00000000 00000000 00000000 00000000 00000000 00000000"),
				labelBitmap: ("" +
					"01001001 00000101 00100000 00100010 00100010 11011001 01010010 01001010," +
					"01011010 10100010 10010101 01010110 10101010 10101010 11010101 00010101," +
					"01001010 10010110 11010111 01101010 01101010 10101010 10101101 01001111," +
					"01101010 10101101 10101010 10101011 01110111 01101010 11011111 01011010," +
					"10101011 11011000 00000000 00000000 00000000 00000000 00000000 00000000"),
				ranks:   []int32{0, 21, 52, 87, 126, 135},
				selects: []int32{1, 151, 260},
				labels:  "Aabnrabdehioubdmnrstaiyeiolnronduiaitaigsudseremimocdirieatcemsiiaisiliactoneeoascthesbndiatnneesirenoaeioedneaacarninlcelhnaaeugnsses",
			},
		},
	}

	for i, c := range cases {
		_ = i
		s := NewSet(c.keys)
		got := wantTyp{
			bitmap.Fmt(s.leaves),
			bitmap.Fmt(s.labelBitmap),
			s.ranks,
			s.selects,
			string(s.labels),
		}
		ta.Equal(c.want, got, "%d-th: struct; case: %+v", i+1, c)

		for _, k := range c.keys {
			found := s.Has(k)
			ta.True(found, "search for %v, case: %v", k, c)
		}

		absent := testutil.RandStrSlice(len(c.keys)*2, 0, 10)
		for _, k := range absent {

			found := s.Has(k)

			idx := sort.SearchStrings(c.keys, k)
			has := idx < len(c.keys) && c.keys[idx] == k

			ta.Equal(has, found, "search for: %v, case: %v", k, c)
		}
	}
}

var OutputSet bool

func BenchmarkSet_Has(b *testing.B) {
	fn := "200kweb2"
	keys := testkeys.Load(fn)

	s := NewSet(keys)

	load := zipf.Accesses(2, 1.5, len(keys), b.N, nil)

	b.ResetTimer()

	var v bool
	for i := 0; i < b.N; i++ {
		rst := s.Has(keys[load[i]])
		v = v || rst
	}
	OutputSet = v
}
