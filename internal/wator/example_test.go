package wator

import (
	"fmt"
	"log"
)

func Example() {
	wator := Wator{}
	if err := wator.init(10, 10, 10, 10, 10, 10); err != nil {
		log.Fatalf(err.Error())
	}
	Update()
	State()
	// Output:
	// String of F, S, *
}

func ExampleNewShark() {

	s := NewShark()
	fmt.Println(s.Age())
	// Output:
	// 0
}

func ExampleWator_Update() {

	Update()
	// Output:
	//  F S * * *

}
