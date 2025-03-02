package wator_test

import (
	"testing"

	"internal/wator"
)

func TestInit(t *testing.T) {
	var w wator.Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Errorf("Failed to initialize world: %v", err)
	}
}

// TestInit tests the initialization of the Wator world.
func TestInit(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Errorf("Failed to initialize world: %v", err)
	}

	if len(w.world) != 25 {
		t.Errorf("Expected world size to be 25, got %d", len(w.world))
	}
}

// TestWorldState tests the State method to check if the snapshot is correct.
func TestWorldState(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	state := w.State()
	if len(state) != 25 {
		t.Errorf("Expected world state size to be 25, got %d", len(state))
	}
}

// TestUpdate tests if the world state updates as expected.
func TestUpdate(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	initialState := w.State()
	w.Update()
	updatedState := w.State()

	// Ensure the world has changed after an update
	if len(initialState) != len(updatedState) {
		t.Fatalf("World state size changed unexpectedly")
	}

	// Ensure the creatures have moved or changed position.
	if initialState == updatedState {
		t.Errorf("World state did not change after update")
	}
}

// TestAdjacent checks if the adjacent function works correctly by testing some corner cases.
func TestAdjacent(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	// Check adjacent for a middle element
	adjacents := w.adjacent(12)
	expected := []int{7, 17, 11, 13} // up, down, left, right for index 12
	for i, adj := range adjacents {
		if adj != expected[i] {
			t.Errorf("Expected adjacent[%d] to be %d, got %d", i, expected[i], adj)
		}
	}

	// Check corner case
	adjacents = w.adjacent(0)
	expected = []int{0, 5, 4, 1} // up, down, left, right for index 0
	for i, adj := range adjacents {
		if adj != expected[i] {
			t.Errorf("Expected adjacent[%d] to be %d, got %d", i, expected[i], adj)
		}
	}
}

// TestDirection tests if the Direction method gives the correct direction between two positions.
func TestDirection(t *testing.T) {
	var w Wator

	tests := []struct {
		start, end int
		expected   int
	}{
		{0, 1, MOVE_EAST},
		{1, 0, MOVE_WEST},
		{0, 5, MOVE_SOUTH},
		{5, 0, MOVE_NORTH},
		{2, 2, MOVE_NONE}, // Same position
	}

	for _, test := range tests {
		t.Run("Direction", func(t *testing.T) {
			result := w.Direction(test.start, test.end)
			if result != test.expected {
				t.Errorf("Expected Direction from %d to %d to be %d, got %d", test.start, test.end, test.expected, result)
			}
		})
	}
}

// TestPickPosition tests if the pickPosition method returns a valid position.
func TestPickPosition(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	// Create a list of open positions for testing.
	openTiles := []int{1, 2, 3, 4, 5}
	chosenPos := w.pickPosition(0, openTiles)

	if chosenPos < 1 || chosenPos > 5 {
		t.Errorf("Expected chosen position to be between 1 and 5, got %d", chosenPos)
	}
}

// TestSpawnFish tests the spawning logic for fish.
func TestSpawnFish(t *testing.T) {
	var w Wator
	err := w.Init(5, 5, 5, 5, 3, 2, 3)
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	// Test that fish spawn correctly after a certain number of updates.
	initialFishCount := 0
	for _, tile := range w.world {
		if _, ok := tile.(*fish); ok {
			initialFishCount++
		}
	}

	// Simulate a few updates
	for i := 0; i < 5; i++ {
		w.Update()
	}

	// Check if the number of fish has increased (indicating spawning)
	newFishCount := 0
	for _, tile := range w.world {
		if _, ok := tile.(*fish); ok {
			newFishCount++
		}
	}

	if newFishCount <= initialFishCount {
		t.Errorf("Fish did not spawn after updates. Expected more fish, got %d", newFishCount)
	}
}
