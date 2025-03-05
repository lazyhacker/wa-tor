// Package wator is an implementation of the wa-tor simulation that A.K. Dewdney
// presented in Scientific America in 1984.
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
)

const (
	NONE  = iota // no creature at the position
	FISH         // Represents a fish in Wator.
	SHARK        // Represents a shark in Wator.
)

// WorldState is the state of Wa-tor at a given Chronon.  The index is the world
// position and the value is what is at the position.
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
	MOVE_EAST         // Movement right
	MOVE_WEST         // Movement left
	DEATH             // creature died
	BIRTH             // New spawn
	ATE               // Creature ate
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
type Wator struct {
	world          []worldItem // Game map is a NxM but represented linearly.
	Width, Height  int         // Dimension of the world.
	Chronon        uint        // Age of the world
	fishSpawnRate  int         // Chronon for a fish to spawn a new fish
	sharkSpawnRate int         // Chronon for a shark to spawn a new shark
}

// Init will set up the world and populate the initial set of fish and shark
// at random positions in the world.  fsr and ssr are the rate by which fish
// and sharks will spawn a new born.  health is the number of Chronon before
// a shark dies if it hasn't eaten a fish.
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
		adjacents := w.adjacentList(i)
		newPos := i
		switch c := tile.(type) {
		case *fish:
			var f *fish

			// Return fish movement and if it spawned a new fish.
			newPos, f = w.fishTurn(c, i, adjacents)

			if f != nil {
				// Put the spawn at the new position because it will
				// then get swapped before the turn is completed and
				// end up in the current position.
				w.world[newPos] = f
				w.world[newPos].setLastMove(w.Chronon)
				w.recordChange(&delta, FISH, i, i, BIRTH)
			}
			w.recordChange(&delta, FISH, i, newPos, MOVE)

		case *shark:

			var s *shark
			var alive bool
			alive, newPos, s = w.sharkTurn(c, i, adjacents)
			// If shark doesn't eat, it dies.
			if !alive {
				w.world[i] = nil
				w.recordChange(&delta, SHARK, i, i, DEATH)
				continue
			}

			w.recordChange(&delta, SHARK, i, newPos, MOVE)

			if _, ok := w.world[newPos].(*fish); ok {
				w.world[newPos] = nil
				w.recordChange(&delta, SHARK, i, newPos, ATE)
			}
			if s != nil {
				w.world[newPos] = s
				w.world[newPos].setLastMove(w.Chronon)
				w.recordChange(&delta, SHARK, i, i, BIRTH)
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

// fishTurns handles the action of a fish each turn and returns its new position
// and if it spawned a new fish.
func (w *Wator) fishTurn(fish *fish, pos int, adjacents []int) (int, *fish) {

	newPos := fish.move(pos, w.world, adjacents)
	if fish.spawn() && newPos != pos {
		return newPos, NewFish()
	}

	return newPos, nil

}

// sharkTurn handles a shark's behavior each turn.
func (w *Wator) sharkTurn(shark *shark, pos int, adjacents []int) (bool, int, *shark) {
	// If shark doesn't eat, it dies.
	if shark.starve() == 0 {
		return false, pos, nil
	}

	newPos := shark.move(pos, w.world, adjacents)
	if _, ok := w.world[newPos].(*fish); ok {
		shark.feed()
	}

	// Cannot spawn if no open space.
	if shark.spawn() && newPos != pos {
		return true, newPos, NewShark()
	}

	return true, newPos, nil

}

// recordChange adds a change to the changelog and uses strings for animals and
// actions instead of the numeric values.
func (w *Wator) recordChange(changelog *[]Delta, animal, from, to, action int) {

	if action == MOVE {
		action = w.direction(from, to)
	}

	*changelog = append(*changelog, Delta{
		Object: animal,
		From:   from,
		To:     to,
		Action: action,
	})
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

// adjacentList returns the for adjacent positions a slice.
func (w *Wator) adjacentList(pos int) []int {

	up, down, left, right := w.adjacents(pos)
	return []int{up, down, left, right}
}

// adjacents returns the four positions next to a given point.
func (w *Wator) adjacents(pos int) (up, down, left, right int) {
	//row = pos / w.Width
	//col = pos % w.Width

	totalTiles := w.Width * w.Height

	up = pos - w.Width
	down = pos + w.Width
	left = pos - 1
	right = pos + 1

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
		right -= w.Width
	}

	// Check if it needs to wrap around to the start of the row.
	if left%w.Width == w.Width-1 || left < 0 {
		left += w.Width
	}

	return
}

// pickPosition randomly picks the element from the given slice.
func (w *Wator) pickPosition(curr int, numbers []int) int {

	if len(numbers) == 0 {
		return curr
	}
	return numbers[rand.Intn(len(numbers))]
}

// direction returns the relative direction of the end position to the start
// position.
func (w *Wator) direction(start, end int) int {

	north, south, west, east := w.adjacents(start)

	switch end {
	case north:
		return MOVE_NORTH
	case west:
		return MOVE_WEST
	case east:
		return MOVE_EAST
	case south:
		return MOVE_SOUTH
	}

	return MOVE_NONE
}

// DebugPrint will print out the state of the world.
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
