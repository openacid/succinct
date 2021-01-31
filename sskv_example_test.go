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
