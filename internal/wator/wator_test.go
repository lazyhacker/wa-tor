package wator_test

import (
	"bytes"
	"testing"

	"lazyhacker.dev/wa-tor/internal/wator"
)

func TestInitTooManyCreatures(t *testing.T) {
	var world wator.Wator
	// For a 2x2 world there are 4 cells. Trying to place 3 fish and 2 sharks (5 total) should fail.
	err := world.Init(2, 2, 3, 2, 3, 3, 2)
	if err == nil {
		t.Error("Expected error when placing more creatures than available cells, got nil")
	}
	expected := "Too many creatures to fit on map!"
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}

func TestInitInvalidHealthRate(t *testing.T) {
	var world wator.Wator
	// Here health (5) is greater than the shark spawn rate (3).
	err := world.Init(3, 3, 2, 2, 3, 3, 5)
	if err == nil {
		t.Error("Expected error when shark health is greater than shark spawn rate, got nil")
	}
	expected := "Health meter needs to be less than the Shark spawn rate."
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}

func TestInitSuccess(t *testing.T) {
	var world wator.Wator
	err := world.Init(3, 3, 3, 2, 3, 3, 2)
	if err != nil {
		t.Fatalf("Unexpected error during Init: %v", err)
	}
	state := world.State()
	fishCount, sharkCount := 0, 0
	for _, v := range state {
		switch v {
		case wator.FISH:
			fishCount++
		case wator.SHARK:
			sharkCount++
		}
	}
	if fishCount != 3 {
		t.Errorf("Expected 3 fish, got %d", fishCount)
	}
	if sharkCount != 2 {
		t.Errorf("Expected 2 sharks, got %d", sharkCount)
	}
}

func TestDirection(t *testing.T) {
	var world wator.Wator

	// No movement if start equals end.
	if world.Direction(5, 5) != wator.MOVE_NONE {
		t.Error("Expected MOVE_NONE when start equals end")
	}
	// East: when end equals start + 1.
	if world.Direction(5, 6) != wator.MOVE_EAST {
		t.Error("Expected MOVE_EAST when end equals start+1")
	}
	// West: when end equals start - 1.
	if world.Direction(5, 4) != wator.MOVE_WEST {
		t.Error("Expected MOVE_WEST when end equals start-1")
	}
	// North: when start < end (and not adjacent east).
	if world.Direction(3, 10) != wator.MOVE_NORTH {
		t.Error("Expected MOVE_NORTH when start < end and not adjacent east")
	}
	// South: when start > end (and not adjacent west).
	if world.Direction(10, 3) != wator.MOVE_SOUTH {
		t.Error("Expected MOVE_SOUTH when start > end and not adjacent west")
	}
}

func TestUpdate(t *testing.T) {
	var world wator.Wator
	err := world.Init(3, 3, 3, 2, 3, 3, 2)
	if err != nil {
		t.Fatalf("Unexpected error during Init: %v", err)
	}
	initialState := world.State()
	states := world.Update()

	// Check that Chronon is incremented.
	if world.Chronon != 1 {
		t.Errorf("Expected Chronon to be 1 after Update, got %d", world.Chronon)
	}
	// Verify that the previous state in the result matches the initial state.
	if len(states.Previous) != len(initialState) {
		t.Error("Length of previous state does not match initial state length")
	}
	for i, v := range states.Previous {
		if v != initialState[i] {
			t.Error("Previous state does not match initial state in Update result")
			break
		}
	}
	// Verify that the current state has the same length.
	if len(states.Current) != len(initialState) {
		t.Error("Current state length mismatch in Update result")
	}
	// Although movement is random, we expect some entries in the change log.
	if len(states.ChangeLog) == 0 {
		t.Log("Warning: No changes recorded in the update; this may be due to the randomness of movement.")
	}
}

func TestDebugPrint(t *testing.T) {
	var world wator.Wator
	err := world.Init(3, 3, 2, 2, 3, 3, 2)
	if err != nil {
		t.Fatalf("Unexpected error during Init: %v", err)
	}

	// Redirect output to a buffer if desired. Here we simply call DebugPrint to ensure it doesn't panic.
	var buf bytes.Buffer
	_ = buf // Currently unused; in a more advanced test you might capture and analyze the output.
	world.DebugPrint()
}
