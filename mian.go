package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uuidObj := uuid.New()
	fmt.Println(uuidObj)
	a := 1
	if a == 2 {
		fmt.Println("nihao", a)

	} else if a == 3 {
		fmt.Println("nihao", a)

	} else {
		fmt.Println("nihao", a)
	}
}
