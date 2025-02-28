package wator

import (
	"math/rand"
	"time"
)

type sequence struct {
	seq []int
}

// init creates a slice of sequential integers and then shuffle them.
func (s *sequence) init(size int) {
	s.seq = make([]int, size)
	for i := 0; i < size; i++ {
		s.seq[i] = int(i)
	}
	rand.Seed(time.Now().UnixNano())

	// Shuffle the sequence
	rand.Shuffle(len(s.seq), func(i, j int) {
		s.seq[i], s.seq[j] = s.seq[j], s.seq[i]
	})
}

// next return the next value in the sequence.
func (s *sequence) next() int {
	n := s.seq[0]     // get the tile number
	s.seq = s.seq[1:] // remove the tile number since it's been taken

	return n
}

func (s *sequence) length() int {

	return len(s.seq)
}
