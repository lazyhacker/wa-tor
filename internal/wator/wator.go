// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.
package wator

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	NONE = iota
	FISH
	SHARK
)

var (
	fishSpawnRate  int
	sharkSpawnRate int
	sharkHealth    int
)

type WorldItem interface {
	Age() int
	SetAge(int)
	Spawn() bool
	LastMove() uint
	SetLastMove(uint)
}

type Wator struct {
	world          []WorldItem // Game map is a NxM but represented linearly.
	Width, Height  int
	Chronon        uint
	fishSpawnRate  int
	sharkSpawnRate int
}

func (w *Wator) Init(width, height, numfish, numsharks, fsr, ssr, health int) {

	w.Width = width
	w.Height = height
	fishSpawnRate = fsr
	sharkSpawnRate = ssr
	sharkHealth = health

	mapSize := w.Width * w.Height
	if numfish+numsharks > mapSize {
		log.Fatalf("Too many creatures to fit on map!")
	}

	if health > ssr {
		log.Fatalf("shark spawn rate is faster then health rate so shark will always spawn befor hunger.")
	}

	// Have a sequence of numbers from 0 to mapSize correspond to
	// locations on the world that isn't occupied.
	sequence := Sequence{}
	sequence.Init(mapSize)

	w.world = make([]WorldItem, mapSize)

	// seed fishes on the tile map.
	for i := 0; i < numfish; i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place FISH.")
			break
		}

		p := sequence.Next()
		w.world[p] = NewFish()
	}

	// seed the sharks on the tile map.
	for i := 0; i < numsharks; i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place SHARK.")
			break
		}

		p := sequence.Next()
		w.world[p] = NewShark()
	}
}

// Adjacent returns up, down, left, right tile locations from the position.
func (w *Wator) Adjacent(pos int) []int {

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

func (w *Wator) Update() []int {

	w.Chronon++
	time.Sleep(50 * time.Millisecond)
	for i, tile := range w.world {
		if tile == nil {
			continue
		}
		if tile.LastMove() == w.Chronon {
			continue
		}
		tile.SetLastMove(w.Chronon)
		// Handle movement and creature-specific behaviors.
		adjacent := w.Adjacent(i)
		var openTiles []int
		var newPos int
		switch c := tile.(type) {
		case *Fish:
			// Fish can only move to non-occupied squares.
			for j := 0; j < len(adjacent); j++ {
				if w.world[adjacent[j]] == nil {
					openTiles = append(openTiles, adjacent[j])
				}
			}

			newPos = w.PickPosition(i, openTiles)
			if c.Spawn() {
				w.world[newPos] = NewFish()
				w.world[newPos].SetLastMove(w.Chronon)
			}

		case *Shark:

			//fmt.Printf("(%d) health = %d\n", c.id, c.health)
			//fmt.Printf("%d %d\n", i, c.id)
			// If shark doesn't eat, it dies.
			(*c).health--
			if (*c).health == 0 {
				w.world[i] = nil
				continue
			}

			// Shark cannot move to tiles that have other sharks
			for j := 0; j < len(adjacent); j++ {
				if v, ok := w.world[adjacent[j]].(*Shark); !ok || v == nil {
					openTiles = append(openTiles, adjacent[j])
				}
			}

			newPos = w.PickPosition(i, openTiles)
			if _, ok := w.world[newPos].(*Fish); ok {
				(*c).health = sharkHealth + 1
				w.world[newPos] = nil
			}
			if c.Spawn() {
				w.world[newPos] = NewShark()
				w.world[newPos].SetLastMove(w.Chronon)
			}
		}

		if newPos != i {
			w.world[newPos], w.world[i] = w.world[i], w.world[newPos]
		}

		tile.SetAge(tile.Age() + 1)
	}

	m := make([]int, len(w.world))

	var t int
	for i := 0; i < len(w.world); i++ {
		switch w.world[i].(type) {
		case *Fish:
			t = FISH
		case *Shark:
			t = SHARK
		default:
			t = NONE
		}
		m[i] = t
	}

	return m
}

// PickPosition randomly picks the element from the given slice.
func (w *Wator) PickPosition(curr int, numbers []int) int {

	if len(numbers) == 0 {
		return curr
	}
	rand.Seed(time.Now().UnixNano())
	return numbers[rand.Intn(len(numbers))]
}

func (w *Wator) Debug() {

	for i, _ := range w.world {

		if i%w.Width == 0 {
			fmt.Println()
		}
		switch w.world[i].(type) {
		case *Fish:
			fmt.Print("F")
		case *Shark:
			fmt.Print("S")
		default:
			fmt.Print("*")
		}
	}

	fmt.Println()

}
