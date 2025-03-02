// Package wator is an implementation of the wa-tor simulation a.k. dewdney presented
// in scientific america in 1984.
//
// # Usage:
//
//	world := wator.Wator{}
//	world.Init(...)
//
//	world.Update()
package wator

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	NONE  = iota // no creature at the position
	FISH         // Represents a fish in Wator.
	SHARK        // Represents a shark in Wator.
)

type WorldState []int

// WorldStates contains the positions of every fish and shark on the map.
// The index is the position and the value is NONE, FISH, or SHARK.
type WorldStates struct {
	Previous  WorldState // Position of Fishes/Shark previous chronon.
	Current   WorldState // Position of Fishes/Shark in current chronon.
	ChangeLog []Delta    // List of changes between Chronon.
}

const (
	NO_ACTION  = iota // No action by creature
	MOVE              // Movement
	MOVE_NONE         // No movement
	MOVE_NORTH        // Movement above
	MOVE_SOUTH        // Movement below
	MOVE_EAST         //Movement right
	MOVE_WEST         //Movement left
	DEATH             // creature died
	BIRTH             // new creature born
	ATE               // creature ate
)

// Delta describes the changes of a creature between two Chronon.
type Delta struct {
	Object int // type of creature: FISH, SHARK
	From   int // position in previous Chronon
	To     int // position in current Chronon
	Action int // Action = NO_ACTION, MOVE, DEATH, BIRTH
}

var (
	fishSpawnRate  int
	sharkSpawnRate int
	sharkHealth    int
)

// worlditem is what is at a location on the world map.  This is generally
// a creature or nothing at all.
type worldItem interface {
	age() int
	setAge(int)
	spawn() bool
	lastMove() uint
	setLastMove(uint)
}

// Wator represents the world of Wa-tor, a toroidal (donut-shaped) sea planet
// consisting of fish and sharks.
type Wator struct {
	world          []worldItem // Game map is a NxM but represented linearly.
	Width, Height  int         // Dimension of the world.
	Chronon        uint        // Age of the world
	fishSpawnRate  int         // Chronon for a fish to spawn a new fish
	sharkSpawnRate int         // Chronon for a shark to spawn a new shark
}

// Init will set up the world and populate the initial set of fish and shark
// at random positions in the world.
func (w *Wator) Init(width, height, numfish, numsharks, fsr, ssr, health int) error {

	w.Width = width
	w.Height = height
	fishSpawnRate = fsr
	sharkSpawnRate = ssr
	sharkHealth = health

	mapSize := w.Width * w.Height
	if numfish+numsharks > mapSize {
		return fmt.Errorf("Too many creatures to fit on map!")
	}

	// If sharks spawns faster then health meter drop rate the popluataion
	// will never decrease.
	if health > ssr {
		return fmt.Errorf("Health meter needs to be less than the Shark spawn rate.")
	}

	// Have a sequence of numbers that will get randomnized to determine
	// where to initially seed the world.
	sequence := sequence{}
	sequence.init(mapSize)

	w.world = make([]worldItem, mapSize)

	// seed fishes on the tile map.
	for i := 0; i < numfish; i++ {

		if sequence.length() == 0 {
			log.Println("No more tiles left on map to place FISH.")
			break
		}

		p := sequence.next()
		w.world[p] = NewFish()
	}

	// seed the sharks on the tile map.
	for i := 0; i < numsharks; i++ {

		if sequence.length() == 0 {
			log.Println("No more tiles left on map to place SHARK.")
			break
		}

		p := sequence.next()
		w.world[p] = NewShark()
	}

	return nil
}

// Update advances the world by 1 Chronon.  During each Chronon:
//   - Fish feed on ubiuitous plankton and the sharks feed on the fish.
//   - Fish move randomly to an unoccupied adjacent square.
//   - After a number of chronon, a fish will spawn another fish.
//   - Sharks will move to an adjacent square if there is a fish and eats the
//     fish otherwise it will move to an random adjacent unoccupied square.
//   - Sharks must eat a fish within a number of cycles or it will die.
//   - At a certain age a shark will spawn a new shark.
func (w *Wator) Update() WorldStates {

	prev := w.State()
	w.Chronon++
	var delta []Delta

	for i, tile := range w.world {
		if tile == nil {
			continue
		}
		// If the creature was already moved this cycle then skip it.
		if tile.lastMove() == w.Chronon {
			continue
		}

		// update to indicate that creature had a turn
		tile.setLastMove(w.Chronon)

		// find adjacent positions
		adjacents := w.adjacent(i)
		var openTiles []int
		var newPos int
		switch c := tile.(type) {
		case *fish:
			// Fish can only move to non-occupied squares.
			for j := 0; j < len(adjacents); j++ {
				if w.world[adjacents[j]] == nil {
					openTiles = append(openTiles, adjacents[j])
				}
			}

			newPos = w.pickPosition(i, openTiles)
			// TODO: This can be refactored to use an interface method
			// and not repeat similar behavior with shark.
			/*
				if c.spawn() {
					// Put the spawn at the new position because it will
					// then get swapped before the turn is completed and
					// end up in the current position.
					w.world[newPos] = NewFish()
					w.world[newPos].setLastMove(w.Chronon)
					delta = append(delta, Delta{
						Object: FISH,
						From:   i,
						To:     i,
						Action: BIRTH,
					})
				}
			*/
			delta = append(delta, Delta{
				Object: FISH,
				From:   i,
				To:     newPos,
				Action: w.Direction(i, newPos),
			})

		case *shark:

			// If shark doesn't eat, it dies.
			(*c).health--
			if (*c).health == 0 {
				w.world[i] = nil
				delta = append(delta, Delta{
					Object: SHARK,
					From:   i,
					To:     i,
					Action: DEATH,
				})
				continue
			}

			// Shark cannot move to tiles that have other sharks
			for j := 0; j < len(adjacents); j++ {
				switch w.world[adjacents[j]].(type) {
				case *fish:
					// If there is a fish, go to that position.
					openTiles = nil
					openTiles = append(openTiles, adjacents[j])
					continue
				case nil:
					openTiles = append(openTiles, adjacents[j])
				}
			}

			newPos = w.pickPosition(i, openTiles)
			delta = append(delta, Delta{
				Object: SHARK,
				From:   i,
				To:     newPos,
				Action: w.Direction(i, newPos),
			})
			if _, ok := w.world[newPos].(*fish); ok {
				(*c).health = sharkHealth + 1
				w.world[newPos] = nil
				delta = append(delta, Delta{
					Object: SHARK,
					From:   i,
					To:     newPos,
					Action: ATE,
				})
			}
			if c.spawn() {
				w.world[newPos] = NewShark()
				w.world[newPos].setLastMove(w.Chronon)
				delta = append(delta, Delta{
					Object: SHARK,
					From:   i,
					To:     i,
					Action: BIRTH,
				})
			}
		}

		if newPos != i {
			// Move the creature by swapping its current location with new position
			w.world[newPos], w.world[i] = w.world[i], w.world[newPos]
		}

		tile.setAge(tile.age() + 1)
	}

	current := w.State()

	return WorldStates{
		Previous:  prev,
		Current:   current,
		ChangeLog: delta,
	}
}

// State returns the snapshop of where each fish and shark is at on the map.
func (w *Wator) State() []int {
	wm := make([]int, len(w.world))
	var t int
	for i := 0; i < len(w.world); i++ {
		switch w.world[i].(type) {
		case *fish:
			t = FISH
		case *shark:
			t = SHARK
		default:
			t = NONE
		}
		wm[i] = t
	}
	return wm
}

// adjacent returns up, down, left, right tile locations from the position.
func (w *Wator) adjacent(pos int) []int {

	totalTiles := w.Width * w.Height
	up := pos - w.Width
	down := pos + w.Width
	left := pos - 1
	right := pos + 1

	// Check if needs to loop around to the bottom of the map.
	if up < 0 {
		up += totalTiles
	}

	// Check to see if needs to loop around to the top of the map.
	if down >= totalTiles {
		down -= totalTiles
	}

	// Check if it needs to go to wrap around to the end of the row.
	if (right % w.Width) == 0 {
		right -= w.Height
	}

	// Check if it needs to wrap around to the start of the row.
	if (left % w.Width) < 0 {
		left += w.Height
	}

	return []int{up, down, left, right}
}

// pickPosition randomly picks the element from the given slice.
func (w *Wator) pickPosition(curr int, numbers []int) int {

	if len(numbers) == 0 {
		return curr
	}
	rand.Seed(time.Now().UnixNano())
	return numbers[rand.Intn(len(numbers))]
}

func (w *Wator) DebugPrint() {

	for i, _ := range w.world {

		if i%w.Width == 0 {
			fmt.Println()
		}
		switch w.world[i].(type) {
		case *fish:
			fmt.Print("F")
		case *shark:
			fmt.Print("S")
		default:
			fmt.Print("*")
		}
	}

	fmt.Println()

}

func (w *Wator) Direction(start, end int) int {

	if start == end {
		return MOVE_NONE
	}
	if end == (start + 1) {
		return MOVE_EAST
	}
	if end == (start - 1) {
		return MOVE_WEST
	}
	if start < end {
		return MOVE_NORTH
	}
	return MOVE_SOUTH

}
