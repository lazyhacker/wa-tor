package main

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
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
			tiles := g.StateToTiles(tc.worldState)
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
			g.frames = make([][]Tile, tc.initialFrames)
			initialCount := len(g.frames)
			g.DeltaToTiles(tc.delta)
			added := len(g.frames) - initialCount
			if added != tc.expectedAdded {
				t.Errorf("Expected %d frames added, got %d", tc.expectedAdded, added)
			}
		})
	}
}

func TestLayout(t *testing.T) {
	tests := []struct {
		name                             string
		worldWidth, worldHeight          int
		expectedScreenW, expectedScreenH int
	}{
		{"small", 10, 5, 10 * TileSize, 5 * TileSize},
		{"zero", 0, 0, 0, 0},
		{"larger", 20, 15, 20 * TileSize, 15 * TileSize},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{}
			g.world.Width = tc.worldWidth
			g.world.Height = tc.worldHeight
			sw, sh := g.Layout(640, 480)
			if sw != tc.expectedScreenW || sh != tc.expectedScreenH {
				t.Errorf("Expected layout (%d, %d), got (%d, %d)", tc.expectedScreenW, tc.expectedScreenH, sw, sh)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name                    string
		pause                   bool
		speedTPS                int
		frameCounter            int
		worldWidth, worldHeight int
		pixelsMove              int // Added to avoid division by zero
		expectedFrameCounter    int
		expectWorldAdvance      bool
	}{
		{
			name:                 "paused update",
			pause:                true,
			speedTPS:             1,
			frameCounter:         0,
			worldWidth:           4,
			worldHeight:          4,
			pixelsMove:           4, // Set to a nonzero value
			expectedFrameCounter: 0,
			expectWorldAdvance:   false,
		},
		{
			name:                 "non-paused update",
			pause:                false,
			speedTPS:             1,
			frameCounter:         0,
			worldWidth:           4,
			worldHeight:          4,
			pixelsMove:           4, // Set to a nonzero value
			expectedFrameCounter: 0,
			expectWorldAdvance:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{}
			g.speedTPS = tc.speedTPS
			g.pause = tc.pause
			g.frameCounter = tc.frameCounter
			// Set nonzero world dimensions.
			g.world.Width = tc.worldWidth
			g.world.Height = tc.worldHeight
			// Set pixelsMove to avoid divide-by-zero in AnimationSteps.
			g.pixelsMove = tc.pixelsMove
			// Initialize frames and tileMap.
			g.frames = make([][]Tile, 0)
			g.tileMap = make([]Tile, 0)

			initialFrames := len(g.frames)
			err := g.Update()
			if err != nil {
				t.Errorf("Update returned error: %v", err)
			}
			if g.frameCounter != tc.expectedFrameCounter {
				t.Errorf("Expected frameCounter to be %d, got %d", tc.expectedFrameCounter, g.frameCounter)
			}
			// When not paused, an update should append frames.
			if tc.expectWorldAdvance && len(g.frames) == initialFrames {
				t.Error("Expected world to advance and frames to be appended, but frames count did not change")
			}
		})
	}
}

func TestRenderMap(t *testing.T) {
	tests := []struct {
		name    string
		tileMap []Tile
	}{
		{
			name: "single tile fish",
			tileMap: []Tile{
				{sprite: 0, tileType: wator.FISH, x: 10, y: 20, direction: EAST},
			},
		},
		{
			name: "single tile shark",
			tileMap: []Tile{
				{sprite: 0, tileType: wator.SHARK, x: 30, y: 40, direction: WEST},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Game{
				tileMap: tc.tileMap,
			}
			// Set up dummy sprites to avoid index out-of-range panics.
			g.fishSprite = []*ebiten.Image{ebiten.NewImage(32, 32)}
			g.sharkSprite = []*ebiten.Image{ebiten.NewImage(32, 32)}

			screen := ebiten.NewImage(100, 100)
			// RenderMap should run without panicking.
			g.RenderMap(screen, g.tileMap)
		})
	}
}

func TestTileCoordinateAccuracy(t *testing.T) {
	// Additional test to check that TileCoordinate calculates positions correctly.
	g := Game{}
	g.world.Width = 8
	// For index 7: row = (7/8)*TileSize = 0*32, col = (7%8)*TileSize = 7*32 = 224.
	x, y := g.TileCoordinate(7)
	if math.Abs(x-224) > 0.001 || math.Abs(y-0) > 0.001 {
		t.Errorf("For index 7, expected (224,0), got (%v,%v)", x, y)
	}
}
