package main

import "lab1/binTree"
import "lab1/custom"
import "fmt"
import "math"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *binTree.Tree, ch chan int) {
	_walk(t, ch)
	close(ch)
}

func _walk(t *binTree.Tree, ch chan int) {
	if t != nil {
		_walk(t.Left, ch)
		ch <- t.Value
		_walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *binTree.Tree) bool {
	if binTree.Compare(t1, t2) {
		return true
	}
	return false
}
func main() {
	fmt.Println("--------------------------------------------")
	fmt.Println("Sorting and print the tree")
	fmt.Println("--------------------------------------------")
	ch := make(chan int)
	go Walk(binTree.New(10, 1), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("Comparing two trees. First should match, second should not")
	fmt.Println("--------------------------------------------")
	if Same(binTree.New(10, 1), binTree.New(10, 1)) {
		fmt.Println("Is Same")
	} else {
		fmt.Println("Is not Same")
	}
	if Same(binTree.New(10, 1), binTree.New(10, 2)) {
		fmt.Println("Is Same")
	} else {
		fmt.Println("Is not Same")
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("Running loop 'n' times")
	fmt.Println("--------------------------------------------")
	loop(5)
	fmt.Println("--------------------------------------------")
	fmt.Println("Value of ok")
	fmt.Println("--------------------------------------------")
	mystery()
	fmt.Println("--------------------------------------------")
	fmt.Println("Custom sqrt function with negative numbers")
	fmt.Println("--------------------------------------------")
	res, err := custom.Sqrt(-2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("Custom sqrt function")
	fmt.Println("--------------------------------------------")
	res2, err2 := custom.Sqrt(2)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(res2)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("Math.Sqrt function")
	fmt.Println("--------------------------------------------")
	fmt.Println(math.Sqrt(2))
}

func loop(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("Iteration ", i)
	}
}

func mystery() {
	someMap := make(map[int]string)
	someMap[0] = "String"
	_, ok := someMap[0]
	fmt.Println(ok)
}
