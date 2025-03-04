package wator

import (
	"fmt"
	"log"
)

// Existing Examples
func Example() {
	var w Wator
	// Note: For examples we use valid parameters.
	if err := w.Init(10, 10, 10, 10, 10, 10, 5); err != nil {
		log.Fatal(err.Error())
	}
	w.Update()
	w.State()
}

func ExampleNewShark() {
	s := NewShark()
	fmt.Println(s.age())
	// Output:
	// 0
}

func ExampleWator_Update() {
	var w Wator
	if err := w.Init(5, 5, 5, 5, 3, 3, 2); err != nil {
		log.Fatal(err.Error())
	}
	w.Update()
}
