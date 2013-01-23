package custom

import (
	"fmt"
	"math"
)

type MyError struct {
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("ERROR! %s", e.What)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return -1, &MyError{
			"No negative numbers please",
		}
	} else {
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
	return 0, nil
}
