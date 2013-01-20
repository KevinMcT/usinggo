package main

import (
	"fmt"
	"math"
)

type MyError struct {
	number float64
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Cannot take Sqrt of negative number! Conflict with value %v", e.number)
}

func run(x float64) error {
	return &MyError{x}
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return -1, &MyError{x}
	}
	z := float64(2.)
	s := float64(0)
	for {
		z = z - (z*z-x)/(2*z)
		if math.Abs(s-z) < 1e-15 {
			break
		}
		s = z
	}
	return s, nil
}

func main() {
	v, err := Sqrt(2)
	fmt.Println(v)
	if err != nil {
		fmt.Println(err)
	}
	v2, err2 := Sqrt(-2)
	fmt.Println(v2)
	if err2 != nil {
		fmt.Println(err2)
	}
}
