package gvalue

import (
	"fmt"
	"testing"
)

func TestOnce(t *testing.T) {
	func1 := Once(func() int {
		fmt.Println("func1")
		return 1
	})

	func1()
	func1()

}
