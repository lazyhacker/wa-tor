package wator

import (
	"math/rand"
	"time"
)

type Sequence struct {
	sequence []int
}

// Init creates a slice of sequential integers and then shuffle them.
func (s *Sequence) Init(size int) {
	s.sequence = make([]int, size)
	for i := 0; i < size; i++ {
		s.sequence[i] = int(i)
	}
	rand.Seed(time.Now().UnixNano())

	// Shuffle the sequence
	rand.Shuffle(len(s.sequence), func(i, j int) {
		s.sequence[i], s.sequence[j] = s.sequence[j], s.sequence[i]
	})
}

// Next return the next value in the sequence.
func (s *Sequence) Next() int {
	n := s.sequence[0]          // get the tile number
	s.sequence = s.sequence[1:] // remove the tile number since it's been taken

	return n
}

func (s *Sequence) Length() int {

	return len(s.sequence)
}
