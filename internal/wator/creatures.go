package wator

import (
	"math/rand"
)

var sharkID int

type creature struct {
	chronon int
	turn    uint
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

type shark struct {
	health int
	creature
	id int
}

func NewShark() *shark {
	id := sharkID
	sharkID++
	return &shark{
		sharkHealth,
		creature{},
		id,
	}
}

func (s *shark) spawn() bool {

	if s.chronon%sharkSpawnRate == 0 && s.chronon > 0 {
		return true
	}
	return false
}

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

func (s *shark) starve() int {

	s.health--

	return s.health
}

func (s *shark) feed() {
	s.health = sharkHealth
}

type fish struct {
	creature
}

func NewFish() *fish {
	return &fish{
		creature{},
	}
}

func (f *fish) spawn() bool {

	if f.chronon%fishSpawnRate == 0 && f.chronon > 0 {
		return true
	}
	return false
}

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
