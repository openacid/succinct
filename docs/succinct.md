# succinct
--
    import "github.com/openacid/succinct"

Package succinct provides several succinct data types.

## Usage

#### type Set

```go
type Set struct {
}
```

Set is a succinct, sorted and static string set impl with compacted trie as
storage. The space cost is about half lower than the original data.


Implementation

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

In this way every node has a `0` pointing to it(except the root node) and has a
corresponding `1` for it:

                                   .-----.
                           .--.    | .---|-.
                           |.-|--. | | .-|-|.
                           || ↓  ↓ | | | ↓ ↓↓
    labels(ignore space):  ab bx u c y v d øøø
    label bitmap:          0010010101010101111
    node-id:               0  1  2 3 4 5 6 789
                              || | ↑ ↑ ↑ |   ↑
                              || `-|-|-' `---'
                              |`---|-'
                              `----'

To walk from a parent node along a label to a child node, count the number of
`0` upto the bit the label position, then find where the the corresponding `1`
is:

    childNodeId = select1(rank0(i))

In our impl, it is:

    nodeId = countZeros(ss.labelBitmap, ss.ranks, bmIdx+1)
    bmIdx = selectIthOne(ss.labelBitmap, ss.ranks, ss.selects, nodeId-1) + 1

Finally leaf nodes are indicated by another bitmap `leaves`, in which a `1` at
i-th bit indicates the i-th node is a leaf:

    leaves: 0001001111

#### func  NewSet

```go
func NewSet(keys []string) *Set
```
NewSet creates a new *Set struct, from a slice of sorted strings.

#### func (*Set) Has

```go
func (ss *Set) Has(key string) bool
```
Has query for a key and return whether it presents in the Set.

#### func (*Set) String

```go
func (ss *Set) String() string
```
