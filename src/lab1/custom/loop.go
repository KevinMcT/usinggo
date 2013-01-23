package custom

import (
	"fmt"
)

func Loop(n int) string {
	j := 0
	for i := 0; i < n; i++ {
		j++
	}
	return fmt.Sprintf("Number of iterations: %v", j)
}
