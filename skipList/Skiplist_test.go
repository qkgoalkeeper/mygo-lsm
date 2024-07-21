package skipList

import (
	"fmt"
	"testing"
)

func Test_SkipList_Insert(t *testing.T) {
	tree := NewSkipList(24, 0.25)

	tree.Set("a", []byte{1, 2, 3})

	tree.Set("a", []byte{2, 3, 4})

	tree.Set("b", []byte{1, 2, 3})
	tree.Set("c", []byte{1, 2, 3})

	data, _ := tree.Search("a")
	fmt.Println(data)

	data, _ = tree.Search("b")
	fmt.Println(data)
	fmt.Println(tree.GetCount())

	fmt.Println(tree.GetValues())

	tree.Delete("b")

	fmt.Println(tree.GetCount())

	fmt.Println(tree.GetValues())

}
