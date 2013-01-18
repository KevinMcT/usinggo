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
func Same(t1, t2 *BinTree.Tree) bool

func main() {
	ch := make(chan int)
	go Walk(BinTree.New(10,1), ch)
	for i := range ch {
		fmt.Println(i)
	}
}
