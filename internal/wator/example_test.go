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

// TestStateCount checks that after initialization, the state slice contains the expected
// number of fish and sharks.
func TestStateCount(t *testing.T) {
	tests := []struct {
		name               string
		width, height      int
		numFish, numSharks int
	}{
		{"5x5 world", 5, 5, 5, 5},
		{"3x3 world", 3, 3, 2, 2},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			var w Wator
			if err := w.Init(tc.width, tc.height, tc.numFish, tc.numSharks, 3, 3, 2); err != nil {
				t.Fatalf("[%d] Init error: %v", i, err)
			}
			state := w.State()
			fishCount, sharkCount := 0, 0
			for _, s := range state {
				switch s {
				case FISH:
					fishCount++
				case SHARK:
					sharkCount++
				}
			}
			if fishCount != tc.numFish || sharkCount != tc.numSharks {
				t.Errorf("[%d] Expected %d fish and %d sharks, got %d fish and %d sharks", i, tc.numFish, tc.numSharks, fishCount, sharkCount)
			}
		})
	}
}

// TestUpdateEffects verifies that Update increases the world's Chronon and returns
// state arrays with the proper length.
func TestUpdateEffects(t *testing.T) {
	tests := []struct {
		name               string
		width, height      int
		numFish, numSharks int
		fsr, ssr, health   int
	}{
		{"5x5 world update", 5, 5, 3, 3, 3, 3, 2},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			var w Wator
			if err := w.Init(tc.width, tc.height, tc.numFish, tc.numSharks, tc.fsr, tc.ssr, tc.health); err != nil {
				t.Fatalf("[%d] Init error: %v", i, err)
			}
			initialChronon := w.Chronon
			ws := w.Update()
			if w.Chronon != initialChronon+1 {
				t.Errorf("[%d] Expected Chronon to increase by 1, got %d", i, w.Chronon-initialChronon)
			}
			if len(ws.Previous) != tc.width*tc.height || len(ws.Current) != tc.width*tc.height {
				t.Errorf("[%d] State arrays length mismatch, expected %d, got %d and %d", i, tc.width*tc.height, len(ws.Previous), len(ws.Current))
			}
		})
	}
}

// TestDirection uses table-driven tests to verify the Direction method.
func TestDirection(t *testing.T) {
	tests := []struct {
		name           string
		start, end     int
		expectedAction int
	}{
		{"no move", 0, 0, MOVE_NONE},
		{"east", 0, 1, MOVE_EAST},
		{"west", 1, 0, MOVE_WEST},
		{"north", 0, 10, MOVE_NORTH},
		{"south", 10, 0, MOVE_SOUTH},
	}

	var w Wator
	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			got := w.Direction(tc.start, tc.end)
			if got != tc.expectedAction {
				t.Errorf("[%d] Direction(%d, %d) = %d, expected %d", i, tc.start, tc.end, got, tc.expectedAction)
			}
		})
	}
}

//
// Unit Tests for Unexported Methods
//

// TestAdjacent tests the unexported adjacent method.
func TestAdjacent(t *testing.T) {
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
			got := w.adjacent(tc.pos)
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

//
// Unit Tests for Creature Spawn Functions
//

// TestFishSpawn verifies that the fish spawn function behaves correctly.
// It sets a known fishSpawnRate, updates the fish's age, and checks if spawn returns the expected value.
func TestFishSpawn(t *testing.T) {
	// Set fishSpawnRate for testing.
	fishSpawnRate = 3
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"age 0", 0, false},
		{"age 1", 1, false},
		{"age 2", 2, false},
		{"age 3", 3, true},
		{"age 4", 4, false},
		{"age 6", 6, true},
		{"age 7", 7, false},
		{"age 9", 9, true},
	}

	for i, tc := range tests {
		f := NewFish()
		f.setAge(tc.age)
		t.Run(fmt.Sprintf("[%d] Fish age %d", i, tc.age), func(t *testing.T) {
			if got := f.spawn(); got != tc.expected {
				t.Errorf("[%d] Fish spawn() with age %d returned %v, expected %v", i, tc.age, got, tc.expected)
			}
		})
	}
}

// TestSharkSpawn verifies that the shark spawn function behaves correctly.
// It sets a known sharkSpawnRate, updates the shark's age, and checks if spawn returns the expected value.
func TestSharkSpawn(t *testing.T) {
	// Set sharkSpawnRate for testing.
	sharkSpawnRate = 4
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"age 0", 0, false},
		{"age 1", 1, false},
		{"age 2", 2, false},
		{"age 3", 3, false},
		{"age 4", 4, true},
		{"age 5", 5, false},
		{"age 8", 8, true},
		{"age 12", 12, true},
	}

	for i, tc := range tests {
		s := NewShark()
		s.setAge(tc.age)
		t.Run(fmt.Sprintf("[%d] Shark age %d", i, tc.age), func(t *testing.T) {
			if got := s.spawn(); got != tc.expected {
				t.Errorf("[%d] Shark spawn() with age %d returned %v, expected %v", i, tc.age, got, tc.expected)
			}
		})
	}
}

func TestRecordChange(t *testing.T) {
	cases := []struct {
		name              string
		animal            int
		from              int
		to                int
		action            int
		expectedDirection int
	}{
		{"Fish moves east", FISH, 1, 2, MOVE, MOVE_EAST},
		{"Shark moves west", SHARK, 5, 4, MOVE, MOVE_WEST},
		{"Shark dies", SHARK, 3, 3, DEATH, DEATH},
		{"Fish spawns", FISH, 4, 4, BIRTH, BIRTH},
		{"Shark eats", SHARK, 5, 6, ATE, ATE},
		{"Shark moves north", SHARK, 8, 3, MOVE, MOVE_NORTH},
		{"Shark spawns", SHARK, 9, 9, BIRTH, BIRTH},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var changes []Delta
			w := Wator{}
			w.recordChange(&changes, tc.animal, tc.from, tc.to, tc.action)

			if len(changes) != 1 {
				t.Errorf("Expected 1 change log entry, got %d", len(changes))
			}
			if changes[0].Object != tc.animal || changes[0].From != tc.from || changes[0].To != tc.to || changes[0].Action != tc.expectedDirection {
				t.Errorf("Unexpected change log entry: %+v, expected direction: %d", changes[0], tc.expectedDirection)
			}
		})
	}
}

func TestSharkTurn(t *testing.T) {
	cases := []struct {
		name        string
		initialHP   int
		adjacents   []int
		openTiles   []int
		shouldLive  bool
		shouldMove  bool
		shouldSpawn bool
	}{
		{"Shark dies", 1, []int{}, []int{}, false, false, false},
		{"Shark moves to open tile", 3, []int{4, 5}, []int{5}, true, true, false},
		{"Shark eats a fish", 2, []int{3, 6}, []int{3}, true, true, false},
		{"Shark spawns", 4, []int{7, 8}, []int{8}, true, true, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := Wator{}
			shark := NewShark()
			shark.health = tc.initialHP
			w.world = make([]worldItem, 10)
			for _, pos := range tc.adjacents {
				w.world[pos] = nil // Ensure these positions are open
			}
			if len(tc.openTiles) > 0 {
				w.world[tc.openTiles[0]] = NewFish() // Place a fish if applicable
			}

			alive, newPos, baby := w.sharkTurn(shark, 5, tc.adjacents)

			if alive != tc.shouldLive {
				t.Errorf("Shark should be alive: %v, but got: %v", tc.shouldLive, alive)
			}
			if tc.shouldMove && newPos == 5 {
				t.Errorf("Shark should have moved but stayed at position %d", newPos)
			}
			if tc.shouldSpawn && baby == nil {
				t.Errorf("Shark should have spawned but did not")
			}
			if !tc.shouldSpawn && baby != nil {
				t.Errorf("Shark should not have spawned but did")
			}
		})
	}
}
