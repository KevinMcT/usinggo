package main

import (
	"fmt"
	"lab1/binTree"
	"lab1/custom"
)

func main() {
	for {
		fmt.Println("---------------------------------------------------------------------")
		fmt.Println("| Choose the function you would like to test using numbers 1 - 5    |")
		fmt.Println("| 1 - Sort and print tree                                           |")
		fmt.Println("| 2 - Compare two trees                                             |")
		fmt.Println("| 3 - Loop n times. Enter number to loop                            |")
		fmt.Println("| 4 - Value of the mystery code in lab                              |")
		fmt.Println("| 5 - Square root of float64                                        |")
		fmt.Println("| Ctrl + c to quit                                                  |")
		fmt.Println("---------------------------------------------------------------------")
		var in int
		fmt.Scanf("%d", &in)

		switch in {
		case 1:
			fmt.Println("--------------------------------------------")
			fmt.Println("Sorting and print the tree")
			fmt.Println("--------------------------------------------")
			ch := make(chan int)
			go Walk(binTree.New(10, 1), ch)
			for i := range ch {
				fmt.Println(i)
			}
		case 2:
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
		case 3:
			fmt.Println("--------------------------------------------")
			fmt.Println("Enter number of times to loop")
			fmt.Println("--------------------------------------------")
			var loopn int
			fmt.Scanf("%d", &loopn)
			fmt.Println("--------------------------------------------")
			fmt.Println(fmt.Sprintf("Running loop %v times", loopn))
			fmt.Println("--------------------------------------------")
			fmt.Println(custom.Loop(loopn))
		case 4:
			fmt.Println("--------------------------------------------")
			fmt.Println("Value of ok")
			fmt.Println("--------------------------------------------")
			fmt.Println(custom.Mystery())
		case 5:
			fmt.Println("--------------------------------------------")
			fmt.Println("Square root of float. If negative value is chosen, an error will be produced.")
			fmt.Println("Enter number to 'root' now: ")
			var square float64
			fmt.Scanf("%v", &square)
			fmt.Println("--------------------------------------------")
			res, err := custom.Sqrt(square)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}
		}
	}
}

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
