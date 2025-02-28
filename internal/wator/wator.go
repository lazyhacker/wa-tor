// wa-tor is an implementation of the wa-tor simulation a.k. dewdney presented
// in scientific america in 1984.
package wator

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	NONE  = iota // no creature at the position
	FISH         // a wator.Fish
	SHARK        // a wator.Shark
)

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
//
// Example usage:
//
//	world := wator.Wator{}
//	world.Init(...)
//
//	world.Update()
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
func (w *Wator) Update() []int {

	w.Chronon++
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
			if c.spawn() {
				// Put the spawn at the new position because it will
				// then get swapped before the turn is completed and
				// end up in the current position.
				w.world[newPos] = NewFish()
				w.world[newPos].setLastMove(w.Chronon)
			}

		case *shark:

			// If shark doesn't eat, it dies.
			(*c).health--
			if (*c).health == 0 {
				w.world[i] = nil
				continue
			}

			// Shark cannot move to tiles that have other sharks
			for j := 0; j < len(adjacents); j++ {
				if v, ok := w.world[adjacents[j]].(*shark); !ok || v == nil {
					openTiles = append(openTiles, adjacents[j])
				}
			}

			// TODO: instead of randomoly moving, prioritize space with fish
			newPos = w.pickPosition(i, openTiles)
			if _, ok := w.world[newPos].(*fish); ok {
				(*c).health = sharkHealth + 1
				w.world[newPos] = nil
			}
			if c.spawn() {
				w.world[newPos] = NewShark()
				w.world[newPos].setLastMove(w.Chronon)
			}
		}

		if newPos != i {
			// Move the creature by swapping its current location with new position
			w.world[newPos], w.world[i] = w.world[i], w.world[newPos]
		}

		tile.setAge(tile.age() + 1)
	}

	m := make([]int, len(w.world))

	//TODO: Break its out into its own method
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
		m[i] = t
	}

	return m
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
