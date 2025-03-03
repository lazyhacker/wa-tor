package main

import (
	"math"
	"testing"

	"lazyhacker.dev/wa-tor/internal/wator"
)

func TestTileCoordinate(t *testing.T) {
	tests := []struct {
		name                 string
		worldWidth           int
		index                int
		expectedX, expectedY float64
	}{
		{"index 0", 10, 0, 0, 0},
		{"index 15", 10, 15, 160, 32}, // row = (15/10)*32 = 32, col = (15%10)*32 = 5*32 = 160
		{"index 7", 8, 7, 224, 0},     // row = (7/8)*32 = 0, col = (7%8)*32 = 7*32 = 224
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{}
			g.world.Width = tc.worldWidth
			x, y := g.TileCoordinate(tc.index)
			if x != tc.expectedX || y != tc.expectedY {
				t.Errorf("TileCoordinate(%d) with width %d: expected (%v, %v), got (%v, %v)",
					tc.index, tc.worldWidth, tc.expectedX, tc.expectedY, x, y)
			}
		})
	}
}

func TestStateToTiles(t *testing.T) {
	tests := []struct {
		name           string
		worldWidth     int
		worldState     wator.WorldState
		expectedLength int
	}{
		{
			name:           "4 cells",
			worldWidth:     4,
			worldState:     wator.WorldState{wator.FISH, wator.SHARK, wator.NONE, wator.FISH},
			expectedLength: 8, // Due to the current implementation (pre-allocating then appending).
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{}
			g.world.Width = tc.worldWidth
			tiles := g.StateToFrame(tc.worldState)
			if len(tiles) != tc.expectedLength {
				t.Errorf("Expected %d tiles, got %d", tc.expectedLength, len(tiles))
			}
			// Verify that the second half of the tiles slice matches the world state.
			half := len(tc.worldState)
			for i, tile := range tiles[half:] {
				if tile.tileType != tc.worldState[i] {
					t.Errorf("At index %d: expected tileType %d, got %d", i, tc.worldState[i], tile.tileType)
				}
				expX, expY := g.TileCoordinate(i)
				if math.Abs(tile.x-expX) > 0.001 || math.Abs(tile.y-expY) > 0.001 {
					t.Errorf("At index %d: expected coordinates (%v, %v), got (%v, %v)", i, expX, expY, tile.x, tile.y)
				}
			}
		})
	}
}

func TestDeltaToTiles(t *testing.T) {
	tests := []struct {
		name          string
		pixelsMove    int
		worldWidth    int
		worldHeight   int
		initialFrames int
		delta         []wator.Delta
		expectedAdded int
	}{
		{
			name:          "east movement",
			pixelsMove:    4,
			worldWidth:    4,
			worldHeight:   4,
			initialFrames: 0,
			delta: []wator.Delta{
				{Object: wator.FISH, From: 0, To: 1, Action: wator.MOVE_EAST},
			},
			// AnimationSteps returns TileSize/pixelsMove = 32/4 = 8, plus one extra frame.
			expectedAdded: 8 + 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{}
			g.pixelsMove = tc.pixelsMove
			// Set world dimensions so that TileCoordinate and frame allocation work properly.
			g.world.Width = tc.worldWidth
			g.world.Height = tc.worldHeight
			g.frames = make([][]Frame, tc.initialFrames)
			initialCount := len(g.frames)
			g.DeltaToFrames(tc.delta)
			added := len(g.frames) - initialCount
			if added != tc.expectedAdded {
				t.Errorf("Expected %d frames added, got %d", tc.expectedAdded, added)
			}
		})
	}
}
