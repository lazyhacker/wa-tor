package wator

import (
	"math/rand"
)

type creature struct {
	chronon int  // age of the creature.
	turn    uint // the chronon when it last moved.
}

func (c *creature) setAge(a int) {
	c.chronon = a
}

func (c *creature) age() int {
	return c.chronon
}

func (c *creature) lastMove() uint {
	return c.turn
}

func (c *creature) setLastMove(t uint) {

	c.turn = t
}

// shark is the predetor on Wa-tor and feeds off fish.
type shark struct {
	health int
	creature
}

// NewShark returns a new instance of a Shark.
func NewShark() *shark {
	return &shark{
		sharkHealth,
		creature{},
	}
}

// spawn returns whether it should spawn a new shark.
func (s *shark) spawn() bool {

	if s.chronon%sharkSpawnRate == 0 && s.chronon > 0 {
		return true
	}
	return false
}

// move determines how a shark moves.
func (s *shark) move(pos int, world []worldItem, adjacents []int) int {

	var openTiles []int
	// Shark cannot move to tiles that have other sharks
	for j := 0; j < len(adjacents); j++ {
		switch world[adjacents[j]].(type) {
		case *fish:
			// If there is a fish, go to that position.
			openTiles = nil
			openTiles = append(openTiles, adjacents[j])
			break
		case nil:
			openTiles = append(openTiles, adjacents[j])
		}
	}

	return pickPosition(pos, openTiles)

}

// starve adjusts the health of a shark when it doesn't eat.
func (s *shark) starve() int {

	s.health--

	return s.health
}

// feed adjusts the shark's health when it eats a fish.
func (s *shark) feed() {
	s.health = sharkHealth
}

// fish is a creature of Wa-tor who eats the planktons in the water.  They
// provide food to sharks.
type fish struct {
	creature
}

// NewFish returns a new instance of a fish.
func NewFish() *fish {
	return &fish{
		creature{},
	}
}

// spawn determines whether a new fish should spawn.
func (f *fish) spawn() bool {

	if f.chronon%fishSpawnRate == 0 && f.chronon > 0 {
		return true
	}
	return false
}

// move handles the fish's movement.
func (f *fish) move(pos int, world []worldItem, adjacents []int) int {

	var openTiles []int
	// Fish can only move to non-occupied squares.
	for j := 0; j < len(adjacents); j++ {
		if world[adjacents[j]] == nil {
			openTiles = append(openTiles, adjacents[j])
		}
	}

	return pickPosition(pos, openTiles)
}

// pickPosition randomly picks the element from the given slice.
func pickPosition(curr int, numbers []int) int {

	if len(numbers) == 0 {
		return curr
	}
	return numbers[rand.Intn(len(numbers))]
}
