package main

import "usinggo/BinTree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *BinTree.Tree, ch chan int) {
	_walk(t, ch)
	close(ch)
}

func _walk(t *BinTree.Tree, ch chan int) {
	if t != nil {
		_walk(t.Left, ch)
		ch <- t.Value
		_walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *BinTree.Tree) bool {
	if BinTree.Compare(t1, t2) {
		return true
	}
	return false
}
func main() {
	ch := make(chan int)
	go Walk(BinTree.New(10, 1), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("Comparing two trees. First should match, second should not")
	fmt.Println("--------------------------------------------")
	if Same(BinTree.New(10, 1), BinTree.New(10, 1)) {
		fmt.Println("Is Same")
	} else {
		fmt.Println("Is not Same")
	}
	if Same(BinTree.New(10, 1), BinTree.New(10, 2)) {
		fmt.Println("Is Same")
	} else {
		fmt.Println("Is not Same")
	}
}
