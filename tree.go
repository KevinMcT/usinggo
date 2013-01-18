package main

import "tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int){
	_walk(t, ch)
	close(ch)
}

func _walk(t *tree.Tree, ch chan int) {
	if t != nil {
		_walk(t.Left, ch)
		fmt.Print("A")
		fmt.Println(t.Value)
		ch <- t.Value
		_walk(t.Right, ch)
		fmt.Print("B")
		fmt.Println(t.Value)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool

func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := range ch {
		fmt.Println(i)
	}
}
