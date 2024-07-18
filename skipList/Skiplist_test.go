package skipList

import (
	"fmt"
	"math"
	"testing"
)

func Test_SkipList_Insert(t *testing.T) {
	tree := NewSkipList(18, 1/math.E)

	tree.Set("a", []byte{1, 2, 3})

	tree.Set("a", []byte{2, 3, 4})

	tree.Set("b", []byte{1, 2, 3})
	tree.Set("c", []byte{1, 2, 3})

	data := tree.Search("a")
	fmt.Println(data.KV)

	data = tree.Search("b")
	fmt.Println(data.KV)
}
