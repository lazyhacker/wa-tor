package wator

import (
	"fmt"
	"testing"
)

//
// Existing Examples
//
/*
func Example() {
	var w Wator
	// Note: For examples we use valid parameters.
	if err := w.Init(10, 10, 10, 10, 10, 10, 5); err != nil {
		log.Fatalf(err.Error())
	}
	w.Update()
	w.State()
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
	var w Wator
	if err := w.Init(5, 5, 5, 5, 3, 3, 2); err != nil {
		log.Fatalf(err.Error())
	}
	w.Update()
	// Output:
	// F S * * *
}
*/
//
// Unit Tests (Table-Driven)
//

// TestInit tests the Init method for various configurations, including error conditions.
func TestInit(t *testing.T) {
	tests := []struct {
		name               string
		width, height      int
		numFish, numSharks int
		fsr, ssr, health   int
		wantErr            bool
	}{
		{"Too many creatures", 5, 5, 20, 10, 3, 3, 2, true},
		{"Health error", 5, 5, 5, 5, 5, 5, 10, true},
		{"Valid init", 5, 5, 5, 5, 3, 3, 2, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			var w Wator
			err := w.Init(tc.width, tc.height, tc.numFish, tc.numSharks, tc.fsr, tc.ssr, tc.health)
			if (err != nil) != tc.wantErr {
				t.Errorf("[%d] Init() error = %v; expected error: %v", i, err, tc.wantErr)
			}
		})
	}
}

//
// Unit Tests for Unexported Methods
//

// TestAdjacent tests the unexported adjacent method.
func TestAdjacentList(t *testing.T) {
	fmt.Println("TestAdjacent")
	tests := []struct {
		name          string
		width, height int
		pos           int
		expected      []int
	}{
		{"corner 0 in 5x6", 5, 6, 0, []int{25, 5, 4, 1}},
		{"middle 12 in 5x6", 5, 6, 12, []int{7, 17, 11, 13}},
		{"corner 4 in 5x6", 5, 6, 4, []int{29, 9, 3, 0}},
		{"corner 29 in 5x6", 5, 6, 29, []int{24, 4, 28, 25}},
		{"corner 25 in 5x6", 5, 6, 25, []int{20, 0, 29, 26}},
		{"left-middle 10 in 5x6", 5, 6, 10, []int{5, 15, 14, 11}},
		{"right-middle 14 in 5x6", 5, 6, 14, []int{9, 19, 13, 10}},
		{"top-middle 2 in 5x5", 5, 6, 2, []int{27, 7, 1, 3}},
		{"bottom-middle 27 in 5x6", 5, 6, 27, []int{22, 2, 26, 28}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			w := Wator{Width: tc.width, Height: tc.height}
			got := w.adjacentList(tc.pos)
			if len(got) != len(tc.expected) {
				t.Errorf("[%d] adjacent(%d) = %v, expected %v", i, tc.pos, got, tc.expected)
				return
			}
			for j, v := range got {
				if v != tc.expected[j] {
					t.Errorf("[%d] adjacent(%d)[%d] = %d, expected %d", i, tc.pos, j, v, tc.expected[j])
				}
			}
		})
	}
}

// TestPickPosition tests the unexported pickPosition method.
func TestPickPosition(t *testing.T) {
	tests := []struct {
		name        string
		curr        int
		numbers     []int
		expectedSet []int // if numbers is non-empty, result must be in expectedSet; if empty, result should equal curr.
	}{
		{"empty numbers", 5, []int{}, nil},
		{"single element", 5, []int{7}, []int{7}},
		{"multiple elements", 3, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
	}

	var w Wator
	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			result := w.pickPosition(tc.curr, tc.numbers)
			if len(tc.numbers) == 0 {
				if result != tc.curr {
					t.Errorf("[%d] Expected %d when numbers is empty, got %d", i, tc.curr, result)
				}
			} else {
				valid := false
				for _, v := range tc.expectedSet {
					if result == v {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("[%d] pickPosition(%d, %v) returned %d, which is not in expected set %v", i, tc.curr, tc.numbers, result, tc.expectedSet)
				}
			}
		})
	}
}
